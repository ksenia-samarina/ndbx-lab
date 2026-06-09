package login

type Domain struct {
	sessionStorage authStorage
	loginStorage   loginStorage
}

func NewDomain(sessionStorage authStorage, loginStorage loginStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		loginStorage:   loginStorage,
	}
}
