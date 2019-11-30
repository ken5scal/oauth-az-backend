package domain

import (
	"net/url"
	"testing"
)

var exampleRedirectUri = "https://client.example.com/cb"
var exampleCode = "SplxlOBeZQQYbYS6WxSbIA"
var examplegrantType = "authorization_code"
var exampleClientId = "s6BhdRkqt3"

func TestInvalidTokenRequest(t *testing.T) {
	u, _ := url.Parse(exampleRedirectUri)
	builder := tokenBuilder{grantType: examplegrantType, code: exampleCode, clientId: exampleClientId, redirectUri: u}

	assertError := func(t *testing.T, b tokenBuilder) {
		err := b.Verify()
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}

		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}
	}

	t.Run("request with empty grantType", func(t *testing.T) {
		builder.grantType = ""
		assertError(t, builder)
	})

	t.Run("request with empty code", func(t *testing.T) {
		builder.code = ""
		assertError(t, builder)
	})

	t.Run("request with empty redirectUri", func(t *testing.T) {
		builder.redirectUri = nil
		assertError(t, builder)
	})

	t.Run("request with empty clientId", func(t *testing.T) {
		builder.clientId = ""
		assertError(t, builder)
	})
}
