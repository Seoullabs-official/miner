package block

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
)

type Block struct {
	Timestamp       int64
	Hash            HexBytes
	PrevHash        HexBytes
	MainBlockHeight int
	MainBlockHash   HexBytes
	Nonce           HexBytes
	Height          int64
	Difficulty      *big.Int
	Miner           HexBytes
	Validator       HexBytes
	ValidatorList   []HexBytes
}
type HexBytes []byte

func (h HexBytes) String() string {
	return string(h) // UTF-8 문자열로 변환
}

// Additional helper functions if needed
func (h HexBytes) ToHex() string {
	return hex.EncodeToString(h)
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
		return err
	}

	// 유효성 검사: 16진수 형식만 허용
	if len(hexStr)%2 != 0 || !isHexString(hexStr) {
		return fmt.Errorf("invalid hex string: %s", hexStr)
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	*h = bytes
	return nil
}

// 헬퍼 함수: 문자열이 유효한 16진수인지 검사
func isHexString(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

func convertToHexBytes(strings []string) []HexBytes {
	hexBytesList := make([]HexBytes, len(strings))
	for i, str := range strings {
		hexBytesList[i] = HexBytes([]byte(str))
	}
	return hexBytesList
}
