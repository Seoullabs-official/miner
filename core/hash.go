package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Seoullabs-official/miner/api/work"
	"github.com/Seoullabs-official/miner/block"
)

func DecodeToString(hash string) string {
	decoded, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Printf("Error decoding Base64: %v", err)
		return ""
	}
	return hex.EncodeToString(decoded)
}

func FindNonceByReturnForHash(nonce work.HexBytes, timestamp int64) []byte {

	data := InitData(nonce, timestamp)
	hash := sha256.Sum256(data)

	return hash[:]
}

func CalculateHashLimit(b *block.Block) (string, error) {
	// 문자열을 big.Int로 변환
	diff := b.Difficulty

	a := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	// 난이도 값 검증 (0 또는 음수 불가)
	if diff.Cmp(big.NewInt(1)) < 0 {
		return "", fmt.Errorf("invalid diff value")
	}

	result := new(big.Int).Div(a, diff)
	hexResult := result.Text(16)
	paddedHexResult := fmt.Sprintf("%064s", hexResult)

	return paddedHexResult, nil
}
func InitData(nonce []byte, timestamp int64) []byte {

	data := bytes.Join(
		[][]byte{
			ToHex(timestamp),
			nonce},
		[]byte{},
	)
	return data
}
func GetHash(nonce work.HexBytes) []byte {
	data := InitData(nonce, time.Now().Unix()) // timestamp hashlimit겨ㄹ과 받을때의 기준으로 표시해서 넣어주기 지금은 임시
	hash := sha256.Sum256(data)
	return hash[:]
}
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
