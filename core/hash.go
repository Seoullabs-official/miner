package core

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Seoullabs-official/miner/api/work"
)

func GenerateRandomNonce() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// SHA-256 해시를 계산
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}
func ToJSONString(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("input is nil")
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal JSON for value: %v, error: %v", v, err)
		return "", err
	}
	return string(jsonBytes), nil
}
func ComputeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input)) // 한번에 해시 계산
	return hex.EncodeToString(hash[:])   // 해시를 헥스 문자열로 변환
}

func CalculateHash(block work.WorkResponse, nonce string) string {
	blockInfo, err := ToJSONString(block)
	if err != nil {
		fmt.Println("Error converting block to JSON string:", err)
		log.Panic(err)
	}
	combinedString := string(block.PrevHash) + blockInfo + nonce
	sha256Hash := ComputeSHA256(combinedString)
	return sha256Hash
}

func FindNonceByReturnForHash(nonce work.HexBytes, timestamp int64) []byte {

	data := InitData(nonce, timestamp)
	hash := sha256.Sum256(data)

	return hash[:]
}

func CalculateHashLimit(difficulty string) (string, error) {
	// 문자열을 big.Int로 변환
	diff := new(big.Int)
	if _, ok := diff.SetString(difficulty, 10); !ok {
		return "", errors.New("invalid difficulty format")
	}

	// 2^256 계산
	maxHash := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	// 난이도 값 검증 (0 또는 음수 불가)
	if diff.Cmp(big.NewInt(1)) < 0 {
		return "", errors.New("difficulty must be greater than 0")
	}

	// 2^256 / difficulty 계산
	limit := new(big.Int).Div(maxHash, diff)

	// 결과를 16진수 문자열로 변환하고 64자리로 패딩
	hexResult := limit.Text(16)
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
