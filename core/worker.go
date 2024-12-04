package core

import (
	"context"
	"encoding/hex"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/Seoullabs-official/miner/api/work"
)

type MiningResult struct {
	Nonce     string
	Timestamp int64
	Hash      string
	Height    int64
	Validator string
	Miner     string
	PrevHash  string
}

func StartWorkers(ctx context.Context, hashLimit string, curBlock work.WorkResponse, loopCount *uint64, cancelFunc context.CancelFunc) ([]byte, int64) {
	numThreads := runtime.NumCPU()
	results := make(chan struct {
		Nonce     []byte
		Timestamp int64
	}, numThreads)
	done := make(chan struct{})
	var once sync.Once

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
						timestamp := time.Now().Unix()

						// 결과 채널로 전송
						results <- struct {
							Nonce     []byte
							Timestamp int64
						}{
							Nonce:     nonceBytes,
							Timestamp: timestamp,
						}

						once.Do(func() { close(done) })
						return
					}
				}
			}
		}(i)
	}

	// 결과 수신 또는 타임아웃
	select {
	case result := <-results:
		return result.Nonce, result.Timestamp
	case <-ctx.Done():
		log.Println("Context timed out before finding a valid hash")
		return nil, 0
	}
}
