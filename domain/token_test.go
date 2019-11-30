package domain

import (
	"net/url"
	"testing"
	"time"
)

var exampleCode = "SplxlOBeZQQYbYS6WxSbIA"
var examplegrantType = "authorization_code"
var exampleClientId = "s6BhdRkqt3"
var exampleUrl, _ = url.Parse("https://client.example.com/cb")
var b = tokenBuilder{}
var builder = tokenBuilder{
	grantType:   examplegrantType,
	code:        exampleCode,
	redirectUri: exampleUrl,
	authzInfo: &AuthorizationInfo{
		AuthzCode:      exampleCode,
		RedirectUri:    exampleUrl,
		CodeExpiration: time.Now().Local().Add(time.Minute * time.Duration(codeExpirationDuration))},
}

func TestTokenRequestWithWrongParameters(t *testing.T) {
	t.Run("request with invalid grantType", func(t *testing.T) {
		b = builder
		b.grantType = "wrong" + b.grantType
		assertTokenRequestError(t, b, tokenUnsupportedGrantType)
	})

	t.Run("request with wrong redirectUri", func(t *testing.T) {
		b = builder
		b.redirectUri, _ = url.Parse("https://wrong.example.com/cb")
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})

	t.Run("request with wrong authz code", func(t *testing.T) {
		b = builder
		b.code = "wrong" + b.code
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})

	t.Run("request with expired authz code", func(t *testing.T) {
		b = builder
		b.authzInfo.CodeExpiration = time.Now().Local().Add(time.Minute * time.Duration(-1*(codeExpirationDuration+1)))
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})

	// this is essentially the same test as "request with expired authz code"
	// but this time, the CodeExpiration is default value which is January 1, year 1, 00:00:00 UTC
	t.Run("request with default code expiration", func(t *testing.T) {
		b = builder
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})
}

func TestTokenRequestWithMissingParameters(t *testing.T) {
	wanted := tokenInvalidRequest

	t.Run("request with empty grantType", func(t *testing.T) {
		b = builder
		b.grantType = ""
		assertTokenRequestError(t, b, wanted)
	})

	t.Run("request with empty redirectUri", func(t *testing.T) {
		b = builder
		b.redirectUri = nil
		assertTokenRequestError(t, b, wanted)
	})

	t.Run("request with empty code", func(t *testing.T) {
		b = builder
		b.code = ""
		assertTokenRequestError(t, b, wanted)
	})
}

func TestTokenRequest(t *testing.T) {
	t.Run("token request", func(t *testing.T) {
		b = builder
		assertNoTokenRequestError(t, b)

		token := b.Build()
		if token.accessToken == "" {
			t.Errorf("wanted an access token, but didn't get one")
		}

		if token.tokenType == "" {
			t.Errorf("wanted a token type, but didn't get one")
		}

		if token.expiresIn.Microseconds() == 0 {
			t.Errorf("wanted token lifetime, but didn't get one")
		}

		if token.refreshToken == "" {
			t.Errorf("wanted a refresh token, but didn't get one")
		}

		if token.scope == "" {
			t.Errorf("wanted a scope, but didn't get one")
		}
	})

	t.Run("token request with no redirect uri specified authz request", func(t *testing.T) {
		b = builder
		b.authzInfo.RedirectUri = nil
		assertNoTokenRequestError(t, b)
	})
}

func assertNoTokenRequestError(t *testing.T, b tokenBuilder) {
	if err := b.Verify(); err != nil {
		t.Errorf("wanted no error, but got %v", err.Error())
	}
}

func assertTokenRequestError(t *testing.T, b tokenBuilder, wantedError string) {
	err := b.Verify()
	if err == nil {
		t.Errorf("wanted an error %v, but didn't get one", wantedError)
		return
	}

	if err.Error() != wantedError {
		t.Errorf("wanted an error %v, but got %v", wantedError, err.Error())
	}
}
