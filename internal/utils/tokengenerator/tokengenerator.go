package tokengenerator

import (
	"crypto/rand"
	"encoding/hex"
)

type Token string

func GenerateToken() Token {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return Token(hex.EncodeToString(b))
}
