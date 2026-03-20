package session

import (
	"context"
	"log"
	"time"
)

func (d *Domain) UpdateSession(ctx context.Context, sid *Sid, ttl time.Duration) error {
	err := d.storage.UpdateSession(ctx, sid, ttl)
	if err != nil {
		log.Printf("Error update session: %v", err)
		return err
	}
	return nil
}
