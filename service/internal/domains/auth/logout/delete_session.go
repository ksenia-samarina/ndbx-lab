package logout

import (
	"context"
	"log"
	"samarina/ndbx/internal/model"
)

func (d *Domain) DeleteSession(ctx context.Context, sid model.Sid) error {
	err := d.storage.DeleteSession(ctx, sid)
	if err != nil {
		log.Printf("Error delete user auth: %v", err)
		return err
	}
	return nil
}
