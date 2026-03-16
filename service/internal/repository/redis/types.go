package redis

type SessionValue struct {
	CreatedAt string `redis:"created_at"`
	UpdatedAt string `redis:"updated_at"`
}
