package main

import (
	"io"
	"log"
	"os"

	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/mining"
	logger "github.com/sirupsen/logrus" // logrus를 log로 별칭 지정
)

func main() {
	// 설정 파일 로드
	cfg, err := config.LoadConfig("../config/config.json")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true, // 타임스탬프 포함
		ForceColors:   true, // 컬러 출력 강제
	})

	logFile, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// 다중 출력: 콘솔 + 파일
	multiOutput := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiOutput)

	// 로그 포맷 설정
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true, // 콘솔 출력에 색상 강제 적용
	})

	// 로그 레벨 설정
	logger.SetLevel(logger.InfoLevel)

	// 초기화 로그
	logger.Info("🚀 Starting the mining process: Initialization complete.")
	logger.Infof("🔧 Miner is configured for target: %s", cfg.TargetMiner)
	logger.Infof("🌐 Establishing connection to domain: %s", cfg.Domain)

	// 마이닝 시작
	mining.Start(cfg)
}
