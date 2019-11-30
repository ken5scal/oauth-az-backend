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
	tokenBuilder := tokenBuilder{}

	t.Run("request with empty grantType", func(t *testing.T) {
		tokenBuilder.grantType = examplegrantType
		err := tokenBuilder.Verify()
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}

		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}
	})

	t.Run("request with empty code", func(t *testing.T) {
		tokenBuilder.code = exampleCode
		err := tokenBuilder.Verify()
		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}
	})

	t.Run("request with empty redirectUri", func(t *testing.T) {
		tokenBuilder.redirectUri, _ = url.Parse(exampleRedirectUri)
		err := tokenBuilder.Verify()
		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}
	})

	t.Run("request with empty clientId", func(t *testing.T) {
		tokenBuilder.clientId = exampleClientId
		err := tokenBuilder.Verify()
		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but didn't get one", tokenInvalidRequest)
		}
	})
}
