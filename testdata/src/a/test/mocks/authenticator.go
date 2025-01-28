package mocks

type Authenticator struct{}

func (a *Authenticator) Authenticate(token string) bool {
	return true
}
