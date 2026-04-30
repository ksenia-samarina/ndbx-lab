package users

import (
	"context"
	"log"
	"samarina/ndbx/internal/domains/types"
	"samarina/ndbx/internal/utils"
	"time"
)

func (d *Domain) RegisterUser(ctx context.Context, user *types.User, ttl time.Duration) (*types.Sid, error) {
	// Регистрируем нового пользователя
	hashedPwd, err := utils.HashedPwd(user.Password)
	if err != nil {
		log.Printf("Error hash password: %v", err)
		return nil, err
	}
	newUser := types.NewUser(user.FullName, user.Username, user.Password, hashedPwd)
	err = d.userStorage.RegisterUser(ctx, newUser)
	if err != nil {
		log.Printf("Error register user in db: %v", err)
		return nil, err
	}
	hexString, err := utils.RandHexString()
	if err != nil {
		log.Printf("Error generating session id: %v", err)
		return nil, err
	}
	// Создаем сессию для нового пользователя
	sid := types.NewSid(hexString)
	userID, err := d.userStorage.GetInternalUserID(ctx, user.Username)
	if err != nil {
		log.Printf("Error get internal user id: %v", err)
		return nil, err
	}
	err = d.sessionStorage.CreateUserSession(ctx, userID.Hex(), sid, ttl)
	if err != nil {
		log.Printf("Error create user session: %v", err)
		return nil, err
	}
	return sid, nil
}
