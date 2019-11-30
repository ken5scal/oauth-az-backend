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
	repo domain.TokenRepository
}

func NewTokenHandler(r domain.TokenRepository) *tokenHandler {
	return &tokenHandler{r}
}

func (c *tokenHandler) RequestToken(w http.ResponseWriter, r *http.Request) {
	// TODO parse request and retrieve parameters
	authzInfo := "" //For now
	token, err := c.GenerateToken(authzInfo)

	// TODO Write to w
	fmt.Println(token)
	fmt.Println(err)
	response.ResponseRequestToken(w, c.GetReturnedName())
}

func (c *tokenHandler) GetTokenByID(tokenID string) (*domain.Token, error) {
	token, err := c.repo.GetAccessTokenByID(tokenID)
	return token, err
}

func (c *tokenHandler) GetReturnedName() string {
	return ""
}

func (c *tokenHandler) GenerateToken(authZInfor string) (*domain.Token, error) {
	// Put Business Logic
	// Check Business logics
	// Insert to db

	r, ok := c.repo.(*infrastructure.TokenRepositoryImpl)
	if !ok {
		return nil, errors.New("TokenRepositoryImpl does not implement TokenRepository")
	}
	r.BeginTransaction()

	token := &domain.Token{}
	err := c.repo.Insert(token)

	if err != nil {
		r.Rollback()
		return nil, err
	}
	r.Commit()

	return token, err
}
