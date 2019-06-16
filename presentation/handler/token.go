package handler

import (
	"fmt"
	"github.com/ken5scal/oauth-az/application"
	"net/http"
)

type tokenHandler struct {
	tokenService application.TokenServiceImpl
}

func NewHandler(s application.TokenServiceImpl) *tokenHandler {
	return &tokenHandler{s}
}

func (c *tokenHandler) RequestToken(w http.ResponseWriter, r *http.Request) {
	// TODO parse request and retrieve parameters
	authzInfo := "" //For now
	token, err := c.tokenService.GenerateToken(authzInfo)

	// TODO Write to w
	fmt.Println(token)
	fmt.Println(err)
}
