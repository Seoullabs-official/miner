package mining

import (
	"context"
	"fmt"
	"math/big"
	"time"

	logger "github.com/sirupsen/logrus" // logrus를 log로 별칭 지정

	"github.com/Seoullabs-official/miner/api"
	"github.com/Seoullabs-official/miner/block"
	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/utils"
	"github.com/Seoullabs-official/miner/work"
)

func Initialize(cfg *config.Config, logger *logger.Logger) {
	var loopCount uint64

	// 작업 채널 초기화
	inCommingBlock := make(chan *block.Block)
	api := api.NewAPI(inCommingBlock, logger)
	api.StartServer(cfg.Port)

	go utils.TrackLoopRate(&loopCount, logger)

	miningRun(inCommingBlock, cfg, logger, &loopCount, api)
}

func miningRun(inCommingBlock chan *block.Block, cfg *config.Config, logger *logger.Logger, loopCount *uint64, api *api.API) {
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
				fmt.Println(api.SendUrl)
				api.SubmitResult(api.SendUrl, &expectBlock)
			}
		}
	}
}

func ReturnMappedBlock(workResponse *block.Block) *block.Block {
	return &block.Block{
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
