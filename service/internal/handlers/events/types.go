package events

type Resp struct {
	Message StatusEventCreation `json:"message"`
}

type StatusEventCreation string

const (
	invalidFieldName     StatusEventCreation = "invalid \"%s\" field"
	invalidParameterName StatusEventCreation = "invalid \"%s\" parameter"
	eventAlreadyExists   StatusEventCreation = "event already exists"
)

type EventID struct {
	ID string `json:"id"`
}
