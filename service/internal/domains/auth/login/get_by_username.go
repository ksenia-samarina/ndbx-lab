package login

import (
	"context"
	"log"
	"samarina/ndbx/internal/domains/types"
)

func (d *Domain) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	hashedUser, err := d.loginStorage.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Error get user by username: %v", err)
		return nil, err
	}
	return hashedUser, nil
}
