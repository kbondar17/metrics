package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Hash(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	rawHash := h.Sum(nil)
	return fmt.Sprintf("%x", rawHash)
}

func HashEqual(h1, h2 []byte) bool {
	return hmac.Equal(h1, h2)
}
