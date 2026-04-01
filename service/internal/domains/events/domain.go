package events

type Domain struct {
	sessionStorage sessionStorage
	eventsStorage  eventsStorage
}

func NewDomain(sessionStorage sessionStorage, loginStorage eventsStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		eventsStorage:  loginStorage,
	}
}
