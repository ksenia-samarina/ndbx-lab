package events

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetInternalUserID(ctx context.Context, sid model.Sid) (string, error) {
	return d.sessionStorage.GetInternalUserID(ctx, sid)
}
