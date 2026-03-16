package session

import (
	"context"
	"log"
	"time"
)

func (d *Domain) UpdateSession(ctx context.Context, id Id, expiration time.Duration) error {
	err := d.Storage.UpdateUserSession(ctx, id, expiration)
	if err != nil {
		log.Printf("Error update session: %v", err)
		return err
	}
	return nil
}
