package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GetSum(metrics []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)

	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}
