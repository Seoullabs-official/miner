package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WorkResponse struct {
	// WorkResponse 구조체 정의
}

func GetWork(domain string) (*WorkResponse, error) {
	url := fmt.Sprintf("%s/getwork", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work: %w", err)
	}
	defer resp.Body.Close()

	var work WorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&work); err != nil {
		return nil, fmt.Errorf("failed to decode work response: %w", err)
	}

	return &work, nil
}

func SubmitResult(domain, nonce string, work *WorkResponse) error {
	url := fmt.Sprintf("%s/completework", domain)
	data := map[string]interface{}{
		"nonce": nonce,
		// 추가 데이터
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
