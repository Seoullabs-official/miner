package main

import (
	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/mining"
	"github.com/Seoullabs-official/miner/utils"
	// logrus를 log로 별칭 지정
)

func main() {
	// 설정 파일 로드
	logger := utils.InitLogger("logfile.log")

	// 설정 파일 로드
	cfg, err := config.LoadConfig("../config/config.json")
	if err != nil {
		logger.Fatalf("❌ Error loading config: %v", err)
	}

	// 초기화 로그
	logger.Info("🚀 Starting the mining process: Initialization complete.")
	logger.Infof("🔧 Miner is configured for target: %s", cfg.TargetMiner)
	logger.Infof("🌐 Establishing connection to domain: %s", cfg.Domain)

	// 마이닝 시작
	mining.Start(cfg, logger)
}
