package logout

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (d *Domain) GetSession(ctx context.Context, sid model.Sid) (bool, error) {
	return d.storage.GetSession(ctx, sid)
}
