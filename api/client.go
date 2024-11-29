package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Seoullabs-official/miner/api/work"
	"github.com/Seoullabs-official/miner/core"
)

func GetWork(domain string) (*work.WorkResponse, error) {
	url := fmt.Sprintf("%s/getwork", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work: %w", err)
	}
	defer resp.Body.Close()

	var work work.WorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&work); err != nil {
		return nil, fmt.Errorf("failed to decode work response: %w", err)
	}

	return &work, nil
}

func SubmitResult(domain string, miningResult *core.MiningResult) error {
	url := fmt.Sprintf("%s/completework", domain)
	data := map[string]interface{}{
		"nonce":     work.HexBytes(miningResult.Nonce), // HexBytes로 변환
		"timestamp": miningResult.Timestamp,
		"height":    miningResult.Height,
		"blockhash": work.HexBytes(miningResult.Hash),      // HexBytes로 변환
		"validator": work.HexBytes(miningResult.Validator), // HexBytes로 변환
		"miner":     work.HexBytes(miningResult.Miner),     // HexBytes로 변환
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 디버깅: JSON 출력
	log.Printf("Submitting JSON: %s", string(jsonData))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	// 서버 응답 디버깅
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned non-OK status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
