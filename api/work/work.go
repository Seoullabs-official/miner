package work

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Seoullabs-official/miner/block"
	"github.com/sirupsen/logrus"
)

// type MiningResult2 struct {
// 	Nonce         HexBytes
// 	Timestamp     int64
// 	Hash          HexBytes
// 	Height        int64
// 	Validator     HexBytes
// 	Miner         HexBytes
// 	PrevHash      HexBytes
// 	Difficulty    *big.Int
// 	ValidatorList []HexBytes
// }

// func (mr MiningResult2) String() string {
// 	var lines []string

// 	// 블록의 기본 정보를 추가
// 	lines = append(lines, "----- Block -----")
// 	lines = append(lines, fmt.Sprintf("Height:      %d", mr.Height))
// 	lines = append(lines, fmt.Sprintf("Timestamp:   %d", mr.Timestamp))
// 	lines = append(lines, fmt.Sprintf("Hash:        %x", mr.Hash))
// 	lines = append(lines, fmt.Sprintf("PrevHash:    %x", mr.PrevHash))
// 	lines = append(lines, fmt.Sprintf("Nonce:       %x", mr.Nonce))
// 	lines = append(lines, fmt.Sprintf("Difficulty:  %d", mr.Difficulty))
// 	lines = append(lines, fmt.Sprintf("Miner: 	    %s", mr.Miner))
// 	lines = append(lines, fmt.Sprintf("Validator:   %s", mr.Validator))
// 	// ValidatorList 정보를 추가
// 	lines = append(lines, "ValidatorList:")

// 	if len(mr.ValidatorList) == 0 {
// 		lines = append(lines, "  (none)")
// 	} else {
// 		for i, v := range mr.ValidatorList {
// 			lines = append(lines, fmt.Sprintf("  %d: %s", i+1, v))
// 		}
// 	}

// 	// 모든 정보를 개행 문자로 구분하여 하나의 문자열로 결합
// 	return strings.Join(lines, "\n")
// }

type HexBytes []byte
type WorkResponse struct {
	Timestamp       int64
	Hash            HexBytes
	PrevHash        HexBytes
	MainBlockHeight int
	MainBlockHash   HexBytes
	Nonce           HexBytes
	Height          int64
	Difficulty      *big.Int // big.Int를 JSON 문자열로 표현
	Miner           HexBytes
	Validator       HexBytes
	ValidatorList   []HexBytes
	// ClientAddress   HexBytes   `json:"client_address"`
}
type ProofOfWork struct {
	Block *block.Block
	Nonce string
}

// type ProofOfWork2 struct {
// 	Block *MiningResult2
// 	Nonce string
// }

func NewProof(b *block.Block) *ProofOfWork {
	nonceStr := fmt.Sprintf("%x", b.Nonce)
	pow := &ProofOfWork{Block: b, Nonce: nonceStr}
	return pow
}

// func NewProof2(b *MiningResult2) *ProofOfWork2 {
// 	nonceStr := fmt.Sprintf("%x", b.Nonce)
// 	pow := &ProofOfWork2{Block: b, Nonce: nonceStr}
// 	return pow
// }

func (pow *ProofOfWork) Validate() bool {
	// 원본 블록의 독립적인 복사본 생성
	validationBlock := *pow.Block // 값 복사

	// 검증용 블록에서 Nonce와 Hash 초기화
	validationBlock.Nonce = nil
	validationBlock.Hash = nil

	// 검증용 블록으로 blockRoot 계산
	blockRoot := CalculateHash(validationBlock, pow.Nonce)

	// 해시 제한 값 계산
	hashLimit, err := CalculateHashLimit(&validationBlock)
	if err != nil {
		log.Println("Error calculating hashLimit:", err)
		return false
	}

	// hashLimit이 blockRoot보다 크거나 같으면 유효한 블록
	isValid := hashLimit >= blockRoot
	return isValid
}
func CalculateHashLimit(b *block.Block) (string, error) {
	// 문자열을 big.Int로 변환
	diff := b.Difficulty

	a := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	// 난이도 값 검증 (0 또는 음수 불가)
	if diff.Cmp(big.NewInt(1)) < 0 {
		return "", fmt.Errorf("invalid diff value")
	}

	result := new(big.Int).Div(a, diff)
	hexResult := result.Text(16)
	paddedHexResult := fmt.Sprintf("%064s", hexResult)

	return paddedHexResult, nil
}
func (pow *ProofOfWork) Run(ctx context.Context, curBlock block.Block, loopCount *uint64, cancelFunc context.CancelFunc) []byte {
	numThreads := runtime.NumCPU()
	results := make(chan struct {
		Nonce []byte
	}, numThreads)
	done := make(chan struct{})
	var once sync.Once
	hashLimit, err := CalculateHashLimit(&curBlock)
	if err != nil {
		logrus.Warnf("Error calculating hash limit: %v", err)
	}
	// 고루틴 생성
	for i := 0; i < numThreads; i++ {
		go func(threadID int) {
			for {
				select {
				case <-done:
					return // 작업 완료 신호를 받으면 종료
				default:

					// 랜덤 nonce 생성 및 해시 계산
					nonce := GenerateRandomNonce()
					hash := CalculateHash(curBlock, nonce)

					nonceBytes, err := hex.DecodeString(nonce)
					if err != nil {
						log.Panic(err)
					}

					if hashLimit >= hash {

						// 결과 채널로 전송
						results <- struct {
							Nonce []byte
						}{
							Nonce: nonceBytes,
						}

						once.Do(func() { close(done) })
						return
					}
					atomic.AddUint64(loopCount, 1)

				}
			}
		}(i)
	}

	// 결과 수신 또는 타임아웃
	select {
	case result := <-results:
		return result.Nonce
	case <-ctx.Done():
		log.Println("Context timed out before finding a valid hash")
		return nil
	}
}
func ComputeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input)) // 한번에 해시 계산
	return hex.EncodeToString(hash[:])   // 해시를 헥스 문자열로 변환
}
func ToJSONString(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("input is nil")
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal JSON for value: %v, error: %v", v, err)
		return "", err
	}
	return string(jsonBytes), nil
}
func CalculateHash(block block.Block, nonce string) string {
	blockInfo, err := ToJSONString(block)
	if err != nil {
		fmt.Println("Error converting block to JSON string:", err)
		log.Panic(err)
	}
	prevHash := fmt.Sprintf("%x", block.PrevHash)
	combinedString := prevHash + blockInfo + nonce
	sha256Hash := ComputeSHA256(combinedString)
	return sha256Hash
}

func GenerateRandomNonce() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// SHA-256 해시를 계산
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}
func (h HexBytes) String() string {
	return string(h) // UTF-8 문자열로 변환
}

// MarshalJSON implements the json.Marshaler interface for HexBytes.
func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil || len(h) == 0 {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%x"`, h)), nil
}

func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Validate hex string
	if len(hexStr)%2 != 0 || !isHexString(hexStr) {
		return fmt.Errorf("invalid hex string: %s", hexStr)
	}

	// Decode hex string
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("hex decoding failed: %w", err)
	}
	*h = bytes
	return nil
}

func isHexString(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

type WorkCompleteResponse struct{}
