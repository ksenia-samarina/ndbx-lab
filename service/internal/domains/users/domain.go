package users

type Domain struct {
	sessionStorage sessionStorage
	userStorage    usersStorage
}

func NewDomain(sessionStorage sessionStorage, userStorage usersStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		userStorage:    userStorage,
	}
}
