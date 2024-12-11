package mining

import (
	"context"
	"math/big"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus" // logrus를 log로 별칭 지정

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/api/work"
	"github.com/Seoullabs-official/miner/block"
	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/utils"
)

func Start(cfg *config.Config, logger *logger.Logger) {
	var loopCount uint64

	// 작업 채널 초기화
	inCommingBlock := make(chan *block.Block)
	api := &api.API{InCommingBlock: inCommingBlock, SendUrl: ""}

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
		if err := http.ListenAndServe("0.0.0.0:7822", nil); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()
	logger.Info("Server is running and awaiting connections.")
}

// processMiningWork 작업 처리
func processMiningWork(inCommingBlock chan *block.Block, cfg *config.Config, logger *logger.Logger, loopCount *uint64, api *api.API) {
	for {
		select {
		case block := <-inCommingBlock:
			logger.Info("Received new work.")

			typeToblock := ReturnMappedBlock(block)
			pow := work.NewProof(typeToblock)
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			cancelFunc()

			nonce := pow.Run(ctx, *typeToblock, loopCount, cancelFunc) //nonce
			if nonce == nil {
				logger.Info("No valid hash found within the time limit. Retrying...")
				continue
			} else {
				pow.Block.Nonce = nonce
				hash := pow.FindNonceByReturnForHash()

				expectBlock := FindNonceReturnMappedBlock(pow.Block.Nonce, pow.Block.Timestamp, cfg.TargetMiner, *typeToblock, hash)
				logger.Infof("Successfully created block: hash =%x ,Nonce=%x, Timestamp=%d", hash, pow.Block.Nonce, pow.Block.Timestamp)

				// http://172.30.1.7:8775
				err := api.SubmitResult(api.SendUrl, &expectBlock)
				if err != nil {
					logger.Errorf("Failed to submit result: %v", err)
				} else {
					logger.Info("Result submitted successfully.")
				}
			}
		}
	}
}

func ReturnMappedBlock(workResponse *block.Block) *block.Block {
	block := &block.Block{
		Timestamp:     workResponse.Timestamp,
		Hash:          block.HexBytes(workResponse.Hash),
		PrevHash:      block.HexBytes(workResponse.PrevHash),
		Nonce:         block.HexBytes(workResponse.Nonce),
		Height:        workResponse.Height,
		Difficulty:    new(big.Int).Set(workResponse.Difficulty),
		Miner:         block.HexBytes(workResponse.Miner),
		Validator:     block.HexBytes(workResponse.Validator),
		ValidatorList: convertHexBytesList(workResponse.ValidatorList),
	}
	return block
}

func convertHexBytesList(input []block.HexBytes) []block.HexBytes {
	output := make([]block.HexBytes, len(input))
	for i, v := range input {
		output[i] = block.HexBytes(v)
	}
	return output
}
func FindNonceReturnMappedBlock(findNonce []byte, timestamp int64, targetMiner string, getWork block.Block, hash []byte) block.Block {
	getWork.Hash = hash
	getWork.Timestamp = timestamp
	getWork.Nonce = findNonce
	getWork.Miner = block.HexBytes(targetMiner)

	return getWork
}
