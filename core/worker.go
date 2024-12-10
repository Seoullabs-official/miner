package core

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Seoullabs-official/miner/api/work"
)

type MiningResult struct {
	Nonce         work.HexBytes
	Timestamp     int64
	Hash          work.HexBytes
	Height        int64
	Validator     work.HexBytes
	Miner         work.HexBytes
	PrevHash      work.HexBytes
	Difficulty    *big.Int
	ValidatorList []work.HexBytes
}

func (mr MiningResult) String() string {
	var lines []string

	// 블록의 기본 정보를 추가
	lines = append(lines, "----- Block -----")
	lines = append(lines, fmt.Sprintf("Height:      %d", mr.Height))
	lines = append(lines, fmt.Sprintf("Timestamp:   %d", mr.Timestamp))
	lines = append(lines, fmt.Sprintf("Hash:        %x", mr.Hash))
	lines = append(lines, fmt.Sprintf("PrevHash:    %x", mr.PrevHash))
	lines = append(lines, fmt.Sprintf("Nonce:       %x", mr.Nonce))
	lines = append(lines, fmt.Sprintf("Difficulty:  %d", mr.Difficulty))
	lines = append(lines, fmt.Sprintf("Miner: 	    %s", mr.Miner))
	lines = append(lines, fmt.Sprintf("Validator:   %s", mr.Validator))
	// ValidatorList 정보를 추가
	lines = append(lines, "ValidatorList:")

	if len(mr.ValidatorList) == 0 {
		lines = append(lines, "  (none)")
	} else {
		for i, v := range mr.ValidatorList {
			lines = append(lines, fmt.Sprintf("  %d: %s", i+1, v))
		}
	}

	// 모든 정보를 개행 문자로 구분하여 하나의 문자열로 결합
	return strings.Join(lines, "\n")
}
