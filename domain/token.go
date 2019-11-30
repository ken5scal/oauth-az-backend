package domain

import (
	"github.com/pkg/errors"
	"net/url"
	"time"
)

const (
	tokenInvalidRequest       = "invalid_request"
	tokenInvalidClient        = "invalid_client"
	tokenInvalidGrant         = "invalid_grant"
	tokenUnauthorizedClient   = "unauthorized_client"
	tokenUnsupportedGrantType = "unsupported_grant_type"
	tokenInvalidScope         = "invalid_scope"
)

// domain logic goes here
type tokenBuilder struct {
	grantType   string
	code        string
	redirectUri *url.URL
	authzInfo   *AuthorizationInfo
}

func NewTokenBuilder(grantType, code string, redirectUri *url.URL, authz *AuthorizationInfo) *tokenBuilder {
	return &tokenBuilder{grantType, code, redirectUri, authz}
}

func (builder *tokenBuilder) Verify() error {
	// Check Logic: https://tools.ietf.org/html/rfc6749#section-4.1.3
	// Error Types: https://tools.ietf.org/html/rfc6749#section-5.2
	if builder.grantType == "" || builder.code == "" || builder.authzInfo.CodeExpiration.IsZero() {
		return errors.New(tokenInvalidRequest)
	}

	if builder.grantType != "authorization_code" {
		return errors.New(tokenUnsupportedGrantType)
	}

	if builder.authzInfo.RedirectUri != nil {
		if builder.redirectUri == nil {
			return errors.New(tokenInvalidRequest)
		}

		if builder.redirectUri.String() != builder.authzInfo.RedirectUri.String() {
			return errors.New(tokenInvalidGrant)
		}
	}

	if builder.code != builder.authzInfo.AuthzCode {
		return errors.New(tokenInvalidGrant)
	}

	if !time.Now().Local().Before(builder.authzInfo.CodeExpiration) {
		return errors.New(tokenInvalidGrant)
	}

	return nil
}

func (builder *tokenBuilder) Build() *Token {
	return &Token{}
}

type Token struct {
	accessToken  string
	tokenType    string
	expiresIn    time.Duration
	refreshToken string
	scope        string
}

func NewToken() *Token {
	return &Token{}
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
