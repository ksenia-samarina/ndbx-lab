package logout

import (
	"context"
	"log"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) DeleteSession(ctx context.Context, sid *types.Sid) error {
	err := d.storage.DeleteSession(ctx, sid)
	if err != nil {
		log.Printf("Error delete user session: %v", err)
		return err
	}
	return nil
}
