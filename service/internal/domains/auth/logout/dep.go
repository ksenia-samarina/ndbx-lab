package logout

import (
	"context"
	"samarina/ndbx/internal/model"
)

type storage interface {
	DeleteSession(ctx context.Context, sid model.Sid) error
}
