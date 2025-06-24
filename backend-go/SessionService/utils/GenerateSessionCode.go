package utils

import (
	"math/rand"
	"time"
)

func GenerateSessionCode(codeLength int) string {
	codeChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = codeChars[rand.Intn(len(codeChars))]
	}
	code := string(b)
	return code
}
