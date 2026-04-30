package events

type Domain struct {
	sessionStorage sessionStorage
	eventsStorage  eventsStorage
}

func NewDomain(sessionStorage sessionStorage, eventsStorage eventsStorage) *Domain {
	return &Domain{
		sessionStorage: sessionStorage,
		eventsStorage:  eventsStorage,
	}
}
