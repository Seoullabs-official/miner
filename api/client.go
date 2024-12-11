package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Seoullabs-official/miner/block"
	"github.com/sirupsen/logrus"
)

type API struct {
	InCommingBlock chan *block.Block // 공개 필드로 변경
	SendUrl        string
}

func (api *API) HandleWork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 요청 본문 읽기
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// payload 구조체로 JSON 디코딩
		var payload struct {
			Data    json.RawMessage `json:"data"`
			SendUrl string          `json:"sendUrl"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Failed to decode JSON: %v", err)
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		// WorkResponse 생성 및 데이터 매핑
		var workResponse block.Block
		if err := json.Unmarshal(payload.Data, &workResponse); err != nil {
			log.Printf("Failed to decode WorkResponse data: %v", err)
			http.Error(w, "Bad request: invalid WorkResponse data", http.StatusBadRequest)
			return
		}

		// sendUrl 값을 ClientAddress에 설정
		// workResponse.ClientAddress = ""http://172.30.1.7:8775""

		// ValidatorList 변환
		api.SendUrl = payload.SendUrl

		// 작업 요청을 채널로 전달
		api.InCommingBlock <- &workResponse

		// 성공 응답 반환
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Work received"))
	}
}
func (api *API) SubmitResult(domain string, miningResult *block.Block) error {
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	url := fmt.Sprintf("%s/completework", domain)
	data := map[string]interface{}{
		"nonce":         miningResult.Nonce, // HexBytes로 변환
		"timestamp":     miningResult.Timestamp,
		"height":        miningResult.Height,
		"hash":          miningResult.Hash,      // HexBytes로 변환
		"validator":     miningResult.Validator, // HexBytes로 변환
		"miner":         miningResult.Miner,     // HexBytes로 변환
		"prevHash":      miningResult.PrevHash,
		"validatorList": miningResult.ValidatorList,
		"difficulty":    miningResult.Difficulty,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 디버깅: JSON 출력
	logrus.Infof("Submitting JSON: %s", string(jsonData))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Error("JSON unmarshal error: %v", err) // 추가된 로그

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
