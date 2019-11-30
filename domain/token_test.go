package domain

import (
	"net/url"
	"testing"
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
	authzInfo:   &authorization{code: exampleCode, redirectUri: exampleUrl}}

func TestTokenRequestWithWrongParameters(t *testing.T) {
	t.Run("request with invalid grantType", func(t *testing.T) {
		b = builder
		b.grantType = "wrong" + b.grantType
		assertTokenRequestError(t, b, tokenUnsupportedGrantType)
	})

	t.Run("request with invalid redirectUri", func(t *testing.T) {
		b = builder
		b.redirectUri, _ = url.Parse("https://wrong.example.com/cb")
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})

	t.Run("request with invalid authz code", func(t *testing.T) {
		b = builder
		b.code = "wrong" + b.code
		assertTokenRequestError(t, b, tokenInvalidGrant)
	})
}

func TestTokenRequestWithMissingParameters(t *testing.T) {
	wanted := tokenInvalidRequest

	t.Run("request with empty grantType", func(t *testing.T) {
		b.grantType = ""
		assertTokenRequestError(t, b, wanted)
	})

	t.Run("request with empty code", func(t *testing.T) {
		b = builder
		b.code = ""
		assertTokenRequestError(t, b, wanted)
	})

	t.Run("request with empty redirectUri", func(t *testing.T) {
		b = builder
		b.redirectUri = nil
		assertTokenRequestError(t, b, wanted)
	})
}

func TestTokenRequestWithWithAuthorizationRequestWithoutRedirectUri(t *testing.T) {
	b = builder
	b.authzInfo.redirectUri = nil
	err := b.Verify()
	if err != nil {
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
