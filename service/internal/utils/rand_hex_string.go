package utils

import (
	"crypto/rand"
	"encoding/hex"
)

const hexStringSize = 16

func RandHexString() (string, error) {
	bytes := make([]byte, hexStringSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
