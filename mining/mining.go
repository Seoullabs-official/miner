package mining

import (
	"context"
	"log"
	"time"

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/core"
	utils "github.com/Seoullabs-official/miner/util"
)

func Start(cfg *config.Config) {
	var loopCount uint64
	go utils.TrackLoopRate(&loopCount)

	for {
		// 작업 요청
		work, err := api.GetWork(cfg.Domain)
		if err != nil {
			log.Printf("Failed to fetch work: %v\n", err)
			continue
		}

		// 난이도 계산
		// hashLimit, err := core.CalculateHashLimit(work.Data.ExpectedBlock.Difficulty)
		hashLimit, err := core.CalculateHashLimit("nil")
		if err != nil {
			log.Printf("Error calculating hash limit: %v\n", err)
			continue
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		// apiData, err := core.StartWorkers(ctx, hashLimit, work.Data.ExpectedBlock.PreviousBlockHash, &loopCount, cancelFunc)
		apiData, err := core.StartWorkers(ctx, hashLimit, "nil", &loopCount, cancelFunc)
		cancelFunc() // 타임아웃 후 컨텍스트 취소
		if err != nil || apiData == "" {
			log.Println("No valid hash found within the time limit")
			continue
		}

		// 결과 전송
		err = api.SubmitResult(cfg.Domain, apiData, work)
		if err != nil {
			log.Printf("Failed to submit result: %v\n", err)
		}
	}
}
