package users

type Domain struct {
	sessionStorage sessionStorage
	userStorage    usersStorage
	eventStorage   eventStorage
}

func NewDomain(sessionStorage sessionStorage, userStorage usersStorage, eventStorage eventStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		userStorage:    userStorage,
		eventStorage:   eventStorage,
	}
}
