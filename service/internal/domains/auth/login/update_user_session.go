package login

import (
	"context"
	"log"
	"samarina/ndbx/internal/domains/types"
	"time"
)

func (d *Domain) UpdateUserSession(ctx context.Context, userID string, sid *types.Sid, ttl time.Duration) error {
	err := d.sessionStorage.UpdateUserSession(ctx, userID, sid, ttl)
	if err != nil {
		log.Printf("Error update session: %v", err)
		return err
	}
	return nil
}
