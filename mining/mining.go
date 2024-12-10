package mining

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	// /getwork 핸들러 등록

	inCommingBlock := make(chan *work.WorkResponse) // 버퍼 추가 가능
	api := &api.API{InCommingBlock: inCommingBlock}

	// 루프 속도 추적
	go utils.TrackLoopRate(&loopCount)

	// /getwork 핸들러 등록
	http.HandleFunc("/getwork", api.HandleWork())

	// 서버 실행
	port := cfg.Port
	log.Printf("Starting server on port %s\n", port)
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 작업 처리 루프
	for {
		select {
		case work := <-inCommingBlock:

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
			fmt.Println(nonce, timestamp)

			expectBlock := FindNonceReturnMappedBlock(nonce, timestamp, cfg.TargetMiner, *work)
			// log.Printf("result : %v\n", &expectBlock)

			// // // 결과 전송
			err = api.SubmitResult(string(work.ClientAddress), &expectBlock)
			if err != nil {
				log.Printf("Failed to submit result: %v\n", err)
			}
		}
		// 작업 요청

	}

}

func FindNonceReturnMappedBlock(findNonce []byte, timestamp int64, targetMiner string, getWork work.WorkResponse) core.MiningResult {
	var expectedBlockStruct core.MiningResult

	hash := core.GetHash(findNonce)
	expectedBlockStruct = core.MiningResult{
		Nonce:         findNonce,
		Timestamp:     timestamp,
		PrevHash:      work.HexBytes(getWork.Hash),
		Validator:     getWork.Validator,
		Miner:         work.HexBytes(targetMiner),
		Hash:          hash,
		Difficulty:    getWork.Difficulty,
		Height:        getWork.Height,
		ValidatorList: getWork.ValidatorList,
	}

	return expectedBlockStruct
}

func convertToHexBytes(strings []string) []work.HexBytes {
	hexBytesList := make([]work.HexBytes, len(strings))
	for i, str := range strings {
		hexBytesList[i] = work.HexBytes([]byte(str))
	}
	return hexBytesList
}
