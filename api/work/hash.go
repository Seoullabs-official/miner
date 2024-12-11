package work

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/Seoullabs-official/miner/block"
)

func DecodeToString(hash string) string {
	decoded, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Printf("Error decoding Base64: %v", err)
		return ""
	}
	return hex.EncodeToString(decoded)
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
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
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
func GenerateRandomNonce() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// SHA-256 해시를 계산
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
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
