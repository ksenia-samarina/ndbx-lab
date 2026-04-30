package events

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

func (d *Domain) UpdateUserSession(ctx context.Context, sid model.Sid, ttl time.Duration) error {
	return d.sessionStorage.UpdateSession(ctx, sid, ttl)
}
