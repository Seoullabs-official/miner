package main

import (
	"log"

	"github.com/Seoullabs-official/miner/config"
	"github.com/Seoullabs-official/miner/mining"
)

func main() {
	// 설정 파일 로드
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 마이닝 시작
	mining.Start(cfg)
}
