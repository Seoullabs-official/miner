package mining

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus" // logrus를 log로 별칭 지정

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/api/work"
	"github.com/Seoullabs-official/miner/block"
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
	logger.Info("/getwork handler registered.")

	// 서버 실행
	startServer(cfg.Port, logger)

	// 작업 처리 루프
	processMiningWork(inCommingBlock, cfg, logger, &loopCount, api)
}

func startServer(port string, logger *logger.Logger) {
	go func() {
		logger.Infof("Starting server on port %s...", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()
	logger.Info("Server is running and awaiting connections.")
}

// processMiningWork 작업 처리
func processMiningWork(inCommingBlock chan *work.WorkResponse, cfg *config.Config, logger *logger.Logger, loopCount *uint64, api *api.API) {
	for {
		select {
		case block := <-inCommingBlock:
			logger.Info("Received new work.")

			typeToblock := ReturnMappedBlock(block)
			pow := work.NewProof(typeToblock)
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			cancelFunc()

			nonce := pow.Run(ctx, *typeToblock, loopCount, cancelFunc)
			newHash := block.() // Assuming nonce is []byte or compatible

			// typeToblock.Hash = newHash
			// logger.Infof("Calculated hash limit: %s", hashLimit)

			// nonce, timestamp := core.StartWorkers(ctx, hashLimit, *work, loopCount, cancelFunc)
			// nonce, timestamp := pow.StartWorkers(ctx, hashLimit, *typeToblock, loopCount, cancelFunc)

			if nonce == nil {
				logger.Info("No valid hash found within the time limit. Retrying...")
				continue
			}

			expectBlock := FindNonceReturnMappedBlock(nonce, timestamp, cfg.TargetMiner, *block)
			logger.Infof("Successfully created block: Nonce=%x, Timestamp=%d", nonce, timestamp)
			// pow2 := work.NewProof2(&expectBlock)
			// pow2.Validate()
			if !pow.Validate() {
				fmt.Println("검증실패 nonce")
				return
			}
			// 결과 전송
			err = api.SubmitResult(string("http://172.30.30.15:8775"), &expectBlock)
			if err != nil {
				logger.Errorf("Failed to submit result: %v", err)
			} else {
				logger.Info("Result submitted successfully.")
			}
		}
	}
}

func ReturnMappedBlock(workResponse *work.WorkResponse) *block.Block {
	block := &block.Block{
		Timestamp:       workResponse.Timestamp,
		Hash:            block.HexBytes(workResponse.Hash),
		PrevHash:        block.HexBytes(workResponse.PrevHash),
		MainBlockHeight: workResponse.MainBlockHeight,
		MainBlockHash:   block.HexBytes(workResponse.MainBlockHash),
		Nonce:           block.HexBytes(workResponse.Nonce),
		Height:          workResponse.Height,
		Difficulty:      new(big.Int).Set(workResponse.Difficulty),
		Miner:           block.HexBytes(workResponse.Miner),
		Validator:       block.HexBytes(workResponse.Validator),
		ValidatorList:   convertHexBytesList(workResponse.ValidatorList),
	}
	return block
}
func convertHexBytesList(input []work.HexBytes) []block.HexBytes {
	output := make([]block.HexBytes, len(input))
	for i, v := range input {
		output[i] = block.HexBytes(v)
	}
	return output
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
