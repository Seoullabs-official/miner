package core

import (
	"context"
	"sync"
	"sync/atomic"
)

func StartWorkers(ctx context.Context, hashLimit, previousBlockHash string, loopCount *uint64, cancelFunc context.CancelFunc) (string, error) {
	var apiData string
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ { // 고루틴 개수
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					nonce := GenerateRandomNonce()
					hash := CalculateHash(previousBlockHash, nonce)
					if hash <= hashLimit {
						apiData = nonce
						cancelFunc()
						return
					}
					atomic.AddUint64(loopCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	return apiData, nil
}
