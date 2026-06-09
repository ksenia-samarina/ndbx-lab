package login

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetSession(ctx context.Context, sid model.Sid) (bool, error) {
	return d.sessionStorage.GetSession(ctx, sid)
}
