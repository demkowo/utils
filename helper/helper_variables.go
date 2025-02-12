package helper

var (
	// Var allows setting mock variables.
	Var varInterface = &helperMock{}
)

type varInterface interface {
	get() *helperMock
	SetExpectedPassword(password string)
	SetExpectedError(err map[string]error)
	SetMock(mock map[string]bool)
}

func (v *helperMock) SetExpectedPassword(password string) {
	v.Password = password
}

func (v *helperMock) SetExpectedError(err map[string]error) {
	v.Error = err
}

func (v *helperMock) SetMock(mock map[string]bool) {
	v.IsMock = mock
}

func (v *helperMock) get() *helperMock {
	return v
}
