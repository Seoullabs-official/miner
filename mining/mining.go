package mining

import (
	"context"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus" // logrus를 log로 별칭 지정

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/api/work"
	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/core"
	"github.com/Seoullabs-official/miner/utils"
)

func Start(cfg *config.Config, logger *logger.Logger) {
	var loopCount uint64

	// 작업 채널 초기화
	inCommingBlock := make(chan *work.WorkResponse)
	api := &api.API{InCommingBlock: inCommingBlock}

	// 루프 속도 추적
	go utils.TrackLoopRate(&loopCount, logger)

	// /getwork 핸들러 등록
	http.HandleFunc("/getwork", api.HandleWork())
	logger.Info("🌐 [INFO] /getwork handler registered.")

	// 서버 실행
	startServer(cfg.Port, logger)

	// 작업 처리 루프
	processMiningWork(inCommingBlock, cfg, logger, &loopCount, api)
}

func startServer(port string, logger *logger.Logger) {
	go func() {
		logger.Infof("🚀 [INFO] Starting server on port %s...", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			logger.Fatalf("❌ [ERROR] Server failed to start: %v", err)
		}
	}()
	logger.Info("✅ [INFO] Server is running and awaiting connections.")
}

// processMiningWork 작업 처리
func processMiningWork(inCommingBlock chan *work.WorkResponse, cfg *config.Config, logger *logger.Logger, loopCount *uint64, api *api.API) {
	for {
		select {
		case work := <-inCommingBlock:
			logger.Info("Received new work.")

			hashLimit, err := core.CalculateHashLimit(work.Difficulty)
			if err != nil {
				logger.Warnf("Error calculating hash limit: %v", err)
				continue
			}
			logger.Infof("Calculated hash limit: %s", hashLimit)

			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			nonce, timestamp := core.StartWorkers(ctx, hashLimit, *work, loopCount, cancelFunc)
			cancelFunc()

			if nonce == nil {
				logger.Info("No valid hash found within the time limit. Retrying...")
				continue
			}

			expectBlock := FindNonceReturnMappedBlock(nonce, timestamp, cfg.TargetMiner, *work)
			logger.Infof("Successfully created block: Nonce=%x, Timestamp=%d", nonce, timestamp)

			// 결과 전송
			err = api.SubmitResult(string(work.ClientAddress), &expectBlock)
			if err != nil {
				logger.Errorf("Failed to submit result: %v", err)
			} else {
				logger.Info("Result submitted successfully.")
			}
		}
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
