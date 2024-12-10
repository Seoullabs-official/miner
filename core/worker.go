package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Seoullabs-official/miner/api/work"
)

type MiningResult struct {
	Nonce         work.HexBytes
	Timestamp     int64
	Hash          work.HexBytes
	Height        int64
	Validator     work.HexBytes
	Miner         work.HexBytes
	PrevHash      work.HexBytes
	ValidatorList []work.HexBytes
	Difficulty    *big.Int
}

func (mr MiningResult) String() string {
	var lines []string

	// 블록의 기본 정보를 추가
	lines = append(lines, "----- Block -----")
	lines = append(lines, fmt.Sprintf("Height:      %d", mr.Height))
	lines = append(lines, fmt.Sprintf("Timestamp:   %d", mr.Timestamp))
	lines = append(lines, fmt.Sprintf("Hash:        %x", mr.Hash))
	lines = append(lines, fmt.Sprintf("PrevHash:    %x", mr.PrevHash))
	lines = append(lines, fmt.Sprintf("Nonce:       %x", mr.Nonce))
	lines = append(lines, fmt.Sprintf("Difficulty:  %d", mr.Difficulty))
	lines = append(lines, fmt.Sprintf("Miner: 	    %s", mr.Miner))
	lines = append(lines, fmt.Sprintf("Validator:   %s", mr.Validator))
	// ValidatorList 정보를 추가
	lines = append(lines, "ValidatorList:")

	if len(mr.ValidatorList) == 0 {
		lines = append(lines, "  (none)")
	} else {
		for i, v := range mr.ValidatorList {
			lines = append(lines, fmt.Sprintf("  %d: %s", i+1, v))
		}
	}

	// 모든 정보를 개행 문자로 구분하여 하나의 문자열로 결합
	return strings.Join(lines, "\n")
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
