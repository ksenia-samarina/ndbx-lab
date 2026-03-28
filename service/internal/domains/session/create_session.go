package session

import (
	"context"
	"log"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) CreateSession(ctx context.Context, ttl time.Duration) (*Sid, error) {
	hexString, err := utils.HexString()
	if err != nil {
		log.Printf("Error generating session id: %v", err)
		return nil, err
	}
	sid := NewSid(hexString)
	err = d.storage.CreateSession(ctx, sid, ttl)
	if err != nil {
		log.Printf("Error create session: %v", err)
		return nil, err
	}
	return sid, nil
}
