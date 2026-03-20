package session

import "context"

func (d *Domain) GetSession(ctx context.Context, sid *Sid) (bool, error) {
	return d.storage.GetSession(ctx, sid)
}
