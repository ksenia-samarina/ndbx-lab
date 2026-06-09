package users

import (
	"context"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) RegisterUser(ctx context.Context, user model.User, ttl time.Duration) (model.Sid, error) {
	hashedPwd, err := utils.HashedPwd(user.Password)
	if err != nil {
		return model.Sid{}, err
	}
	newUser := model.NewUser(user.FullName, user.Username, user.Password, hashedPwd)
	err = d.userStorage.RegisterUser(ctx, newUser)
	if err != nil {
		return model.Sid{}, err
	}
	hexString, err := utils.RandHexString()
	if err != nil {
		return model.Sid{}, err
	}
	sid := model.NewSid(hexString)
	userID, err := d.userStorage.GetInternalUserID(ctx, user.Username)
	if err != nil {
		return model.Sid{}, err
	}
	err = d.sessionStorage.CreateUserSession(ctx, userID.Hex(), sid, ttl)
	if err != nil {
		return model.Sid{}, err
	}
	return sid, nil
}
