package logout

import (
	"context"
	"samarina/ndbx/internal/model"
)

type domain interface {
	DeleteSession(ctx context.Context, sid model.Sid) error
	GetSession(ctx context.Context, sid model.Sid) (bool, error)
}
