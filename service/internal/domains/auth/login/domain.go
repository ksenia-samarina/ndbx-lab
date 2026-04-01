package login

type Domain struct {
	sessionStorage sessionStorage
	loginStorage   loginStorage
}

func NewDomain(sessionStorage sessionStorage, loginStorage loginStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		loginStorage:   loginStorage,
	}
}
