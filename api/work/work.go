package work

import "math/big"

type HexBytes []byte
type WorkResponse struct {
	Timestamp       int64    `json:"timestamp"`
	Hash            string   `json:"hash"`
	PrevHash        string   `json:"prevHash"`
	MainBlockHeight int      `json:"mainBlockHeight"`
	MainBlockHash   string   `json:"mainBlockHash"`
	Nonce           string   `json:"nonce"`
	Height          int64    `json:"height"`
	Difficulty      *big.Int `json:"difficulty"` // big.Int를 JSON 문자열로 표현
	Miner           string   `json:"miner"`
	Validator       string   `json:"validator"`
	ValidatorList   []string `json:"validatorList"`
}

type WorkCompleteResponse struct{}
