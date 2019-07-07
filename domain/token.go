package domain

// domain logic goes here
type Token struct {
	authZInfo string
}

func NewToken(authzInfo string) *Token {
	return &Token{authZInfo: authzInfo}
}

// I think they are business logic...
type ReturningToken struct {
}

// I think they are business logic...
func (t *ReturningToken) Name() string {
	return ""
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
	GetAccessTokenByID(tokenID string) (*Token, error)
	Insert(t *Token) error
	Update(t *Token) error
	Delete(t *Token) error
}

type ReturningTokenService interface {
	GetReturnedName() string
}
