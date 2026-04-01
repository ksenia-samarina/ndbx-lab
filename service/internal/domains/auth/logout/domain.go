package logout

type Domain struct {
	storage storage
}

func NewDomain(storage storage) *Domain {
	return &Domain{
		storage: storage,
	}
}
