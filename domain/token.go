package domain

// domain logic goes here

type Token struct {
	authZInfo string
}

type ReturningToken struct {
}

func (t *Token) isRevoked() bool {
	return false
}

func (t *Token) FindClientID() int {
	return 0
}

func ReturnToken(t *Token) *ReturningToken {
	return nil
}

type TokenRepository interface {
	GetByID(tokenID string) (*Token, error)
	Insert(t *Token) error
	Update(t *Token) error
	Delete(t *Token) error
}
