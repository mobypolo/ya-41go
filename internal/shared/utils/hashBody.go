package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashBody(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
