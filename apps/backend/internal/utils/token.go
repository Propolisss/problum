package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken(n int) string {
	token := make([]byte, n)

	rand.Read(token)

	return base64.StdEncoding.EncodeToString(token)
}
