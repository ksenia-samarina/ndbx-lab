package redis

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func (s *Storage) buildReactionsKey(eventTitle string) string {
	hash := md5.Sum([]byte(eventTitle))
	return fmt.Sprintf("event:%s:reactions", hex.EncodeToString(hash[:]))
}

func (s *Storage) buildReviewsKey(eventTitle string) string {
	hasher := md5.New()
	hasher.Write([]byte(eventTitle))
	titleHash := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("event:%s:reviews", titleHash)
}
