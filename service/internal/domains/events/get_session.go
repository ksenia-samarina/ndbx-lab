package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) GetSession(ctx context.Context, sid *types.Sid) (bool, error) {
	return d.sessionStorage.GetSession(ctx, sid)
}
