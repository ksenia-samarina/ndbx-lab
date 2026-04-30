package login

import (
	"context"
	"log"
	"samarina/ndbx/internal/domains/types"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) CreateUserSession(ctx context.Context, userID string, ttl time.Duration) (*types.Sid, error) {
	hexString, err := utils.RandHexString()
	if err != nil {
		log.Printf("Error generating session id: %v", err)
		return nil, err
	}
	sid := types.NewSid(hexString)
	err = d.sessionStorage.CreateUserSession(ctx, userID, sid, ttl)
	if err != nil {
		log.Printf("Error create session: %v", err)
		return nil, err
	}
	return sid, nil
}
