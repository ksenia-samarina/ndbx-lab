package health

type Resp struct {
	Status StatusHealthCheck `json:"status"`
}

type StatusHealthCheck string
