package domain

import (
	"github.com/pkg/errors"
	"net/url"
)

// domain logic goes here
type tokenBuilder struct {
	grantType   string
	code        string
	redirectUri *url.URL
}

func NewTokenBuilder(grantType, code string, redirectUri *url.URL) *tokenBuilder {
	return &tokenBuilder{grantType, code, redirectUri}
}

func (builder *tokenBuilder) Verify() error {
	// https://tools.ietf.org/html/rfc6749#section-4.1.3
	if builder.grantType == "" || builder.code == "" || builder.redirectUri == nil {
		return errors.New(InvalidRequest)
	}

	if builder.grantType != "authorization_code" {
		return errors.New(tokenUnsupportedGrantType)
	}

	return nil
}

const (
	tokenInvalidRequest       = "invalid_request"
	tokenInvalidClient        = "invalid_client"
	tokenInvalidGrant         = "invalid_grant"
	tokenUnauthorizedClient   = "unauthorized_client"
	tokenUnsupportedGrantType = "unsupported_grant_type"
	tokenInvalidScope         = "invalid_scope"
)

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
