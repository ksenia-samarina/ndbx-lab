package reviews

import "time"

type Domain struct {
	reviewStorage reviewStorage
	eventsStorage eventsStorage
	reviewsCache  reviewsCache
	reviewTTL     time.Duration
}

func NewDomain(eventsStorage eventsStorage, reviewStorage reviewStorage, reviewsCache reviewsCache, reviewTTL time.Duration) *Domain {
	return &Domain{
		reviewStorage: reviewStorage,
		eventsStorage: eventsStorage,
		reviewsCache:  reviewsCache,
		reviewTTL:     reviewTTL,
	}
}
