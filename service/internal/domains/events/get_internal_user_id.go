package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) GetInternalUserID(ctx context.Context, sid *types.Sid) (string, error) {
	return d.sessionStorage.GetInternalUserID(ctx, sid)
}
