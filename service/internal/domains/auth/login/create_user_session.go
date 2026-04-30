package login

import (
	"context"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) CreateUserSession(ctx context.Context, userID string, ttl time.Duration) (model.Sid, error) {
	hexString, err := utils.RandHexString()
	if err != nil {
		return model.Sid{}, err
	}
	sid := model.NewSid(hexString)
	err = d.sessionStorage.CreateUserSession(ctx, userID, sid, ttl)
	if err != nil {
		return model.Sid{}, err
	}
	return sid, nil
}
