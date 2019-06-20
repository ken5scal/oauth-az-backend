package handler

// Handles Business logics and Presentation
// `handler` package handles business logic from
// request and response.
// Specifically following s
// * works on  db transaction
// * guarantees data integrity

// IMPORTANT
// put service layer (business logics) in domain layer

import (
	"errors"
	"fmt"
	"github.com/ken5scal/oauth-az/domain"
	"github.com/ken5scal/oauth-az/handler/response"
	"github.com/ken5scal/oauth-az/infrastructure"
	"net/http"
)

type tokenHandler struct {
	tokenService TokenServiceImpl
}

func NewHandler(s TokenServiceImpl) *tokenHandler {
	return &tokenHandler{s}
}

func (c *tokenHandler) RequestToken(w http.ResponseWriter, r *http.Request) {
	// TODO parse request and retrieve parameters
	authzInfo := "" //For now
	token, err := c.tokenService.GenerateToken(authzInfo)

	// TODO Write to w
	fmt.Println(token)
	fmt.Println(err)
	response.ResponseRequestToken(w, c.tokenService.GetReturnedName())
}

type TokenServiceImpl struct {
	repo domain.TokenRepository
}

func NewService(repository domain.TokenRepository) TokenServiceImpl {
	return TokenServiceImpl{
		repo: repository,
	}
}

func (t *TokenServiceImpl) GetTokenByID(tokenID string) (*domain.ReturningToken, error) {
	hoge, err := t.repo.GetAccessTokenByID(tokenID)
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
