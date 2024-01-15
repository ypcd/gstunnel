package gshash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256Hex(data []byte) string {
	s1 := sha256.Sum256(data)
	return hex.EncodeToString(
		s1[:],
	)
}
