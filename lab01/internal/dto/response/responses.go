package response

type HealthCheck struct {
	Status StatusHealthCheck `json:"status"`
}

type StatusHealthCheck string

const OkHealthCheck StatusHealthCheck = "ok"
