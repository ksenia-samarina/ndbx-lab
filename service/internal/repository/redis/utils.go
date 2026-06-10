package redis

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func (s *Storage) buildKey(eventTitle string) string {
	hash := md5.Sum([]byte(eventTitle))
	return fmt.Sprintf("event:%s:reactions", hex.EncodeToString(hash[:]))
}
