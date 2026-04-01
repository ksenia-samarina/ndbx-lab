package login

import (
	"context"
	"log"
)

func (d *Domain) GetInternalUserID(ctx context.Context, username string) (string, error) {
	objectID, err := d.loginStorage.GetInternalUserID(ctx, username)
	if err != nil {
		log.Printf("Error get user by username: %v", err)
		return "", err
	}
	return objectID.String(), nil
}
