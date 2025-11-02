package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomId() int64 {
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	return n.Int64()

}
