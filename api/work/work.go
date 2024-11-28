package work

type WorkResponse struct {
	Timestamp       int64  `json:"timestamp"`
	Hash            string `json:"hash"`
	PrevHash        string `json:"prevHash"`
	MainBlockHeight int    `json:"mainBlockHeight"`
	MainBlockHash   string `json:"mainBlockHash"`
	Nonce           string `json:"nonce"`
	Height          int64  `json:"height"`
	Difficulty      string `json:"difficulty"`
	Miner           string `json:"miner"`
	Validator       string `json:"validator"`
}

type WorkCompleteResponse struct{}
