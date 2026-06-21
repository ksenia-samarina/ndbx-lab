package reactions

import "time"

type Domain struct {
	reactionsStorage reactionsStorage
	reactionsCache   reactionsCache
	eventsStorage    eventsStorage
	likeTTL          time.Duration
}

func NewDomain(storage reactionsStorage, eventsStorage eventsStorage, cache reactionsCache, likeTTL time.Duration) *Domain {
	return &Domain{
		reactionsStorage: storage,
		eventsStorage:    eventsStorage,
		reactionsCache:   cache,
		likeTTL:          likeTTL,
	}
}
