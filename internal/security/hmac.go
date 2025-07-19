package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMAC(secret, data []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}

func HMACString(secret, data []byte) string {
	return hex.EncodeToString(HMAC(secret, data))
}
