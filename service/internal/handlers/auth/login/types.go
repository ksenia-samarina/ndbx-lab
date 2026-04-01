package login

type Resp struct {
	Message StatusUserLogin `json:"message"`
}

type StatusUserLogin string

const (
	userInvalidCredentials StatusUserLogin = "invalid credentials"
	invalidFieldName       StatusUserLogin = "invalid \"%s\" field"
)
