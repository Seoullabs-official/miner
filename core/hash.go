package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateRandomNonce() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}

func CalculateHash(previousHash, nonce string) string {
	data := previousHash + nonce
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func CalculateHashLimit(difficulty string) (string, error) {
	maxValue := new(big.Int)
	maxValue.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)

	diffValue := new(big.Int)
	if _, ok := diffValue.SetString(difficulty, 10); !ok {
		return "", fmt.Errorf("invalid difficulty value")
	}

	result := new(big.Int).Div(maxValue, diffValue)
	return fmt.Sprintf("%064x", result), nil
}
