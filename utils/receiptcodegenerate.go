package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateReceiptNumber() string {
	now := time.Now().Format("20060102150405") // yyyyMMddHHmmss

	// random 4 digit
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))

	return fmt.Sprintf("RC-%s-%04d", now, n.Int64())
}
