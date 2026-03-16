package session

import (
	"context"
	"log"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) CreateSession(ctx context.Context, expiration time.Duration) (Id, error) {
	hexString, err := utils.GenerateHexString()
	if err != nil {
		log.Printf("Error generating session id: %v", err)
		return Id{HexString: ""}, err
	}
	sessionId := Id{HexString: hexString}
	err = d.Storage.CreateUserSession(ctx, sessionId, expiration)
	if err != nil {
		log.Printf("Error update session: %v", err)
		return Id{HexString: ""}, err
	}
	return sessionId, nil
}
