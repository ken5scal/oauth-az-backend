package domain

import (
	"net/url"
	"testing"
)

var exampleCode = "SplxlOBeZQQYbYS6WxSbIA"
var examplegrantType = "authorization_code"
var exampleClientId = "s6BhdRkqt3"
var u, _ = url.Parse("https://client.example.com/cb")

func TestTokenRequestWithWrongParameters(t *testing.T) {
	builder := tokenBuilder{grantType: "not_authorization_code", code: exampleCode, redirectUri: u}
	wanted := tokenUnsupportedGrantType
	err := builder.Verify()
	if err.Error() != tokenUnsupportedGrantType {
		t.Errorf("wanted an error %v, but got %v", wanted, err.Error())
	}
}

func TestTokenRequestWithMissingParameters(t *testing.T) {
	builder := tokenBuilder{grantType: examplegrantType, code: exampleCode, redirectUri: u}
	wanted := tokenInvalidRequest

	assertError := func(t *testing.T, b tokenBuilder) {
		err := b.Verify()
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", wanted)
		}

		if err.Error() != tokenInvalidRequest {
			t.Errorf("wanted an error %v, but got %v", wanted, err.Error())
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
}
