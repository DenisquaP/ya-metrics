package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCryptography(t *testing.T) {
	msg := []byte("hello world")
	key := "hello"

	hash := GetSum(msg, key)

	h := hmac.New(sha256.New, []byte(key))
	h.Write(msg)

	require.Equal(t, hash, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}
