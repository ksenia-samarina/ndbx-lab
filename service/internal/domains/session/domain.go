package session

type Domain struct {
	Storage Storage
}

func New(storage Storage) *Domain {
	return &Domain{
		Storage: storage,
	}
}
