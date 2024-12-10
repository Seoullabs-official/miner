package work

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
)

type HexBytes []byte
type WorkResponse struct {
	Timestamp       int64      `json:"timestamp"`
	Hash            HexBytes   `json:"hash"`
	PrevHash        HexBytes   `json:"prevHash"`
	MainBlockHeight int        `json:"mainBlockHeight"`
	MainBlockHash   HexBytes   `json:"mainBlockHash"`
	Nonce           HexBytes   `json:"nonce"`
	Height          int64      `json:"height"`
	Difficulty      *big.Int   `json:"difficulty"` // big.Int를 JSON 문자열로 표현
	Miner           HexBytes   `json:"miner"`
	Validator       HexBytes   `json:"validator"`
	ValidatorList   []HexBytes `json:"validatorList"`
	ClientAddress   HexBytes   `json:"client_address"`
}

func (h HexBytes) String() string {
	return string(h) // UTF-8 문자열로 변환
}

// MarshalJSON implements the json.Marshaler interface for HexBytes.
func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil || len(h) == 0 {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%x"`, h)), nil
}

func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Validate hex string
	if len(hexStr)%2 != 0 || !isHexString(hexStr) {
		return fmt.Errorf("invalid hex string: %s", hexStr)
	}

	// Decode hex string
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("hex decoding failed: %w", err)
	}
	*h = bytes
	return nil
}

func isHexString(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

type WorkCompleteResponse struct{}
