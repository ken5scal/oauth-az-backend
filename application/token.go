package application

import (
	"errors"
	"github.com/ken5scal/oauth-az/domain"
	"github.com/ken5scal/oauth-az/infrastructure"
)

// this will have business logic
// * works on  db transaction
// * guarantees data integrity
// * presentation uses this

type TokenServiceImpl struct {
	repo domain.TokenRepository
}

func NewService(repository domain.TokenRepository) TokenServiceImpl {
	return TokenServiceImpl{
		repo: repository,
	}
}

func (t *TokenServiceImpl) GetTokenByID(tokenID string) (*domain.ReturningToken, error) {
	hoge, err := t.repo.GetByID(tokenID)
	return domain.ReturnToken(hoge), err
}

func (t *TokenServiceImpl) GetReturnedName() string {
	return ""
}

func (t *TokenServiceImpl) GenerateToken(authZInfor string) (*domain.ReturningToken, error) {
	// Put Business Logic
	// Check Business logics
	// Insert to db

	r, ok := t.repo.(*infrastructure.TokenRepositoryImpl)
	if !ok {
		return nil, errors.New("TokenRepositoryImpl does not implement TokenRepository")
	}
	r.BeginTransaction()

	token := &domain.Token{authZInfor}
	err := t.repo.Insert(token)

	if err != nil {
		r.Rollback()
		return nil, err
	}
	r.Commit()

	return domain.ReturnToken(token), err
}
