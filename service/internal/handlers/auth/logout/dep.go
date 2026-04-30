package logout

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

type domain interface {
	DeleteSession(ctx context.Context, sid *types.Sid) error
}
