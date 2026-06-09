package session

import (
	"context"
	"log"
	"samarina/ndbx/internal/model"
	"time"
)

func (d *Domain) UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error {
	err := d.storage.UpdateSession(ctx, sid, ttl)
	if err != nil {
		log.Printf("Error update auth: %v", err)
		return err
	}
	return nil
}
