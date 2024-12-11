package work

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Seoullabs-official/miner/block"
	"github.com/sirupsen/logrus"
)

type HexBytes block.HexBytes

type ProofOfWork struct {
	Block *block.Block
	Nonce string
}

func NewProof(b *block.Block) *ProofOfWork {
	nonceStr := fmt.Sprintf("%x", b.Nonce)
	pow := &ProofOfWork{Block: b, Nonce: nonceStr}
	return pow
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
					hash := pow.CalculateHash(&curBlock, nonce)

					nonceBytes, err := hex.DecodeString(nonce)
					if err != nil {
						log.Panic(err)
					}

					if hashLimit >= hash {
						pow.Block.Timestamp = time.Now().Unix()
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
	case <-done:
		// results 채널을 닫지 않음, 필요시 go 루틴에서 닫도록 처리
	}

	return nil
}

func (p *ProofOfWork) InitData(nonce []byte) []byte {

	data := bytes.Join(
		[][]byte{
			ToHex(p.Block.Timestamp),
			nonce},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) FindNonceByReturnForHash() []byte {
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)

	return hash[:]
}
func (pow *ProofOfWork) CalculateHash(block *block.Block, nonce string) string {
	blockInfo, err := ToJSONString(block)
	if err != nil {
		fmt.Println("Error converting block to JSON string:", err)
		log.Panic(err)
	}
	// 해시 계산을 위한 결합 문자열 출력
	combinedString := fmt.Sprintf("%s%s%s", fmt.Sprintf("%x", block.PrevHash), blockInfo, nonce)
	// SHA256 해시 계산
	sha256Hash := ComputeSHA256(combinedString)
	return sha256Hash
}
