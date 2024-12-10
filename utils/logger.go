package utils

import (
	"io"

	logger "github.com/sirupsen/logrus"

	"os"
)

func InitLogger(logFilePath string) *logger.Logger {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("❌ Failed to open log file: %v", err)
	}

	// 다중 출력 설정
	multiOutput := io.MultiWriter(os.Stdout, logFile)

	// 로거 생성 및 설정
	log := logger.New()
	log.SetOutput(multiOutput)
	log.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true, // 콘솔 색상 강제 적용
	})
	log.SetLevel(logger.InfoLevel)

	return log
}
