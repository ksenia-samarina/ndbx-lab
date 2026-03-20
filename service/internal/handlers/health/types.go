package health

type Health struct {
	Status StatusHealthCheck `json:"status"`
}

type StatusHealthCheck string
