package mining

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/api/work"
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

		hashLimit, err := core.CalculateHashLimit(work.Difficulty)
		if err != nil {
			log.Printf("Error calculating hash limit: %v\n", err)
			continue
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

		nonce, timestamp := core.StartWorkers(ctx, hashLimit, *work, &loopCount, cancelFunc)
		cancelFunc()
		if nonce == nil {
			log.Println("No valid hash found within the time limit")
			continue
		}

		expectBlock := FindNonceReturnMappedBlock(nonce, timestamp, cfg.TargetMiner, *work)
		log.Printf("result : %v\n", &expectBlock)

		// // // 결과 전송
		err = api.SubmitResult(cfg.Domain, &expectBlock)
		if err != nil {
			log.Printf("Failed to submit result: %v\n", err)
		}
	}

}

func FindNonceReturnMappedBlock(findNonce []byte, timestamp int64, targetMiner string, getWork work.WorkResponse) core.MiningResult {
	var expectedBlockStruct core.MiningResult

	hash := core.GetHash(findNonce)
	expectedBlockStruct = core.MiningResult{
		Nonce:     hex.EncodeToString(findNonce),
		Timestamp: timestamp,
		Height:    getWork.Height + 1,
		PrevHash:  getWork.Hash,
		Validator: getWork.Validator,
		Miner:     targetMiner,
		Hash:      hex.EncodeToString(hash),
	}

	return expectedBlockStruct
}
