package domain

// domain logic goes here

type Token struct {
	authZInfo string
}

type ReturningToken struct {
	name string
}

// I think they are business logic...
func (t *Token) isRevoked() bool {
	return false
}

// I think they are business logic...
func (t *Token) FindClientID() int {
	return 0
}

// I think they are business logic...
func ReturnToken(t *Token) *ReturningToken {
	return nil
}

type TokenRepository interface {
	GetByID(tokenID string) (*Token, error)
	Insert(t *Token) error
	Update(t *Token) error
	Delete(t *Token) error
}
