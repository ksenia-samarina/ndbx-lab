package update_session

import (
	"context"
	"samarina/ndbx/internal/domains/session"
	"time"
)

type Domain interface {
	UpdateSession(ctx context.Context, id session.Id, expiration time.Duration) error
	CreateSession(ctx context.Context, expiration time.Duration) (session.Id, error)
	CheckSessionExist(ctx context.Context, id session.Id) (bool, error)
}
