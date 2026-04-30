package login

import (
	"context"
)

func (d *Domain) GetInternalUserID(ctx context.Context, username string) (string, error) {
	objectID, err := d.loginStorage.GetInternalUserID(ctx, username)
	if err != nil {
		return "", err
	}
	return objectID.String(), nil
}
