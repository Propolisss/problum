package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHMAC(key, message []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}
