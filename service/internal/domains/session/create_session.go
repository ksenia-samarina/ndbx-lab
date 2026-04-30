package session

import (
	"context"
	"log"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) CreateSession(ctx context.Context, ttl time.Duration) (model.Sid, error) {
	hexString, err := utils.RandHexString()
	if err != nil {
		log.Printf("Error generating auth id: %v", err)
		return model.Sid{}, err
	}
	sid := model.NewSid(hexString)
	err = d.storage.CreateSession(ctx, sid, ttl)
	if err != nil {
		log.Printf("Error create auth: %v", err)
		return model.Sid{}, err
	}
	return sid, nil
}
