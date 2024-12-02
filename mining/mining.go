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

		// 난이도 계산 끝나면 그 기준 timestamp ,blockhash , validator , miner 넣어주기
		hashLimit, err := core.CalculateHashLimit(work.Difficulty)
		if err != nil {
			log.Printf("Error calculating hash limit: %v\n", err)
			continue
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

		result, err := core.StartWorkers(ctx, hashLimit, *work, &loopCount, cancelFunc, cfg.TargetMiner, cfg.Validator)
		cancelFunc()
		if err != nil || result.Nonce == "" {
			log.Println("No valid hash found within the time limit")
			continue
		}

		log.Printf("resuult : %v\n", &result)

		// // // 결과 전송
		err = api.SubmitResult(cfg.Domain, &result)
		if err != nil {
			log.Printf("Failed to submit result: %v\n", err)
		}
	}

}
