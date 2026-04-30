package users

type Resp struct {
	Message StatusUserRegistration `json:"message"`
}

type StatusUserRegistration string

const (
	userAlreadyExists StatusUserRegistration = "user already exists"
	invalidFieldName  StatusUserRegistration = "invalid \"%s\" field"
)
