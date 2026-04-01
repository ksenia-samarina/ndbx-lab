package logout

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

type storage interface {
	DeleteSession(ctx context.Context, sid *types.Sid) error
}
