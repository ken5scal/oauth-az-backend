package domain

import (
	"crypto/rand"
	"encoding/base64"
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
	scope       []string
	authzInfo   *AuthorizationInfo
}

func NewTokenBuilder(grantType, code string, redirectUri *url.URL, authz *AuthorizationInfo) *tokenBuilder {
	return &tokenBuilder{
		grantType:   grantType,
		code:        code,
		redirectUri: redirectUri,
		authzInfo:   authz}
}

func (builder *tokenBuilder) Scope(scope []string) *tokenBuilder {
	builder.scope = scope
	return builder
}

func (builder *tokenBuilder) Verify() error {
	// Check Logic: https://tools.ietf.org/html/rfc6749#section-4.1.3
	// Error Types: https://tools.ietf.org/html/rfc6749#section-5.2
	if builder.grantType == "" || builder.code == "" {
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

// check whether scope in authorization request and scope in token response has same value
func (builder *tokenBuilder) hasSameScope() bool {
	if len(builder.scope) != len(builder.authzInfo.Scope) {
		return false
	}

	for _, v := range builder.scope {
		for k, av := range builder.authzInfo.Scope {
			if v == av {
				break
			}

			if len(builder.authzInfo.Scope) == k+1 {
				return false
			}
		}
	}

	return true
}

func (builder *tokenBuilder) Build() *Token {
	// "minimum of 128 bits of entropy where the probability of an attacker guessing the generated token is less than or equal to 2^(-160) as per [RFC6749] section 10.10"
	//  https://bitbucket.org/openid/fapi/pull-requests/45/bring-access-token-requirements-inline/diff
	// calculated by https://8gwifi.org/passwdgen.jsp
	// Don't use symboles, just numbers and letters from 22 ~ 26
	b := make([]byte, 26)
	rand.Read(b)
	accessToken := base64.StdEncoding.EncodeToString(b)
	token := &Token{
		accessToken: accessToken,
		tokenType:   "Bearer", // https://openid.net/specs/openid-connect-core-1_0.html#TokenResponse
		expiresIn:   time.Duration(3600) * time.Second,
	}

	// https://tools.ietf.org/html/rfc6749#section-5.1
	// if not identical to the scope requested by the client in authz request, it's required
	if !builder.hasSameScope() {
		token.scope = builder.scope
	}

	return token
}

type Token struct {
	accessToken  string
	tokenType    string
	expiresIn    time.Duration
	refreshToken string // FAPIでも必須ではない
	scope        []string
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
