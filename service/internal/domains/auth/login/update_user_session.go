package login

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

func (d *Domain) UpdateUserSession(ctx context.Context, userID string, sid model.Sid, ttl time.Duration) error {
	return d.sessionStorage.UpdateUserSession(ctx, userID, sid, ttl)
}
