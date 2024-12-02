package core

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
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
}

func StartWorkers(ctx context.Context, hashLimit string, curBlock work.WorkResponse, loopCount *uint64, cancelFunc context.CancelFunc,
	targetMiner string, targetValidator string) (MiningResult, error) {

	numThreads := runtime.NumCPU()
	results := make(chan MiningResult, numThreads)
	done := make(chan struct{})
	var once sync.Once

	// 고루틴 생성
	for i := 0; i < numThreads; i++ {
		go func(threadID int) {
			for {
				select {
				case <-done:
					return // 다른 고루틴에서 작업 완료 신호를 받으면 종료
				default:
					// 랜덤 nonce 생성 및 해시 계산
					nonce := GenerateRandomNonce()
					hash := CalculateHash(curBlock, nonce)

					if hashLimit >= hash {
						timestamp := time.Now().Unix()

						hashBytes := FindNonceByReturnForHash(work.HexBytes(nonce), timestamp)
						if hashBytes == nil {
							log.Println("Error: FindNonceByReturnForHash returned nil")
							continue
						}
						// 결과 생성 및 전송
						result := MiningResult{
							Nonce:     nonce,
							Timestamp: timestamp,
							Height:    curBlock.Height + 1,
							Hash:      hash,
							Validator: targetValidator,
							Miner:     targetMiner,
						}

						results <- result

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
		return result, nil
	case <-ctx.Done():
		log.Println("Context timed out before finding a valid hash")
		return MiningResult{}, fmt.Errorf("no valid result found within the time limit")
	}
}
