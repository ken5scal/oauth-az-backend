package domain

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var exampleCode = "SplxlOBeZQQYbYS6WxSbIA"
var examplegrantType = "authorization_code"
var exampleClientId = "s6BhdRkqt3"
var exampleUrl, _ = url.Parse("https://client.example.com/cb")
var exampleScope = []string{"openid", "profile", "email", "phone"}
var b = tokenBuilder{}

func initBuilder() tokenBuilder {
	return tokenBuilder{
		grantType:   examplegrantType,
		code:        exampleCode,
		redirectUri: exampleUrl,
		scope:       exampleScope,
		authzInfo: &AuthorizationInfo{
			AuthzCode:      exampleCode,
			RedirectUri:    exampleUrl,
			CodeExpiration: time.Now().Local().Add(time.Minute * time.Duration(codeExpirationDuration)),
			Scope:          exampleScope,
		},
	}
}

func TestTokenRequestWithWrongParameters(t *testing.T) {
	t.Run("request with invalid grantType", func(t *testing.T) {
		b = initBuilder()
		b.grantType = "wrong" + b.grantType
		assertTokenRequestVerifyError(t, b, tokenUnsupportedGrantType)
	})

	t.Run("request with wrong redirectUri", func(t *testing.T) {
		b = initBuilder()
		b.redirectUri, _ = url.Parse("https://wrong.example.com/cb")
		assertTokenRequestVerifyError(t, b, tokenInvalidGrant)
	})

	t.Run("request with wrong authz code", func(t *testing.T) {
		b = initBuilder()
		b.code = "wrong" + b.code
		assertTokenRequestVerifyError(t, b, tokenInvalidGrant)
	})

	t.Run("request with expired authz code", func(t *testing.T) {
		b = initBuilder()
		b.authzInfo.CodeExpiration = time.Now().Local().Add(time.Minute * time.Duration(-1*(codeExpirationDuration+1)))
		assertTokenRequestVerifyError(t, b, tokenInvalidGrant)
	})

	// this is essentially the same test as "request with expired authz code"
	// but this time, the CodeExpiration is default value which is January 1, year 1, 00:00:00 UTC
	t.Run("request with default code expiration", func(t *testing.T) {
		b = initBuilder()
		b.authzInfo.CodeExpiration = time.Time{}
		assertTokenRequestVerifyError(t, b, tokenInvalidGrant)
	})

	t.Run("request with different scope", func(t *testing.T) {
		b = initBuilder()
		b.scope = append(b.authzInfo.Scope, "user")
		fmt.Println(b.authzInfo.CodeExpiration.String())
		assertTokenRequestVerifyError(t, b, tokenInvalidScope)
	})
}

func TestTokenRequestWithMissingParameters(t *testing.T) {
	wanted := tokenInvalidRequest

	t.Run("request with empty grantType", func(t *testing.T) {
		b = initBuilder()
		b.grantType = ""
		assertTokenRequestVerifyError(t, b, wanted)
	})

	t.Run("request with empty redirectUri", func(t *testing.T) {
		b = initBuilder()
		b.redirectUri = nil
		assertTokenRequestVerifyError(t, b, wanted)
	})

	t.Run("request with empty code", func(t *testing.T) {
		b = initBuilder()
		b.code = ""
		assertTokenRequestVerifyError(t, b, wanted)
	})
}

func TestTokenRequest(t *testing.T) {
	t.Run("token request", func(t *testing.T) {
		b = initBuilder()
		assertNoTokenRequestError(t, b)

		token := b.Build()
		if token.accessToken == "" {
			t.Errorf("wanted an access token, but didn't get one")
		}

		if token.tokenType == "" {
			t.Errorf("wanted a token type, but didn't get one")
		}

		if token.expiresIn.Nanoseconds() == 0 {
			t.Errorf("wanted token lifetime, but didn't get one")
		}
		// refresh token can be an optional, even in FAPI security profile
		//if token.refreshToken == "" {
		//	t.Errorf("wanted a refresh token, but didn't get one")
		//}
	})

	t.Run("token request with no redirect uri specified in authz request", func(t *testing.T) {
		b = initBuilder()
		b.authzInfo.RedirectUri = nil
		assertNoTokenRequestError(t, b)
	})

	t.Run("token request with no scope specified in authz request", func(t *testing.T) {
		b = initBuilder()
		assertNoTokenRequestError(t, b)

		newScope := append(b.authzInfo.Scope, "new-scope")
		token := b.Scope(newScope).Build()

		if len(token.scope) == 0 {
			t.Errorf("wanted scope %v, but got none", token.scope)
		}

		if !reflect.DeepEqual(token.scope, newScope) {
			t.Errorf("wanted scope %v, but got %v", newScope, token.scope)
		}
	})
}

func TestScopeDifference(t *testing.T) {
	t.Run("token request with scope in random order", func(t *testing.T) {
		b = initBuilder()
		b.scope = []string{"profile", "email", "phone", "openid"}
		if !b.hasSameScope() {
			t.Errorf("expect to be right")
		}
	})

	t.Run("token request with scope with different length", func(t *testing.T) {
		b = initBuilder()
		b.scope = append(b.authzInfo.Scope, "user")
		if b.hasSameScope() {
			t.Errorf("expect to be wrong")
		}
	})

	t.Run("token request with scope with different values", func(t *testing.T) {
		b = initBuilder()
		b.scope = []string{"profile", "fake-email", "phone", "openid"}
		if b.hasSameScope() {
			t.Errorf("expect to be wrong")
		}
	})
}

func assertNoTokenRequestError(t *testing.T, b tokenBuilder) {
	if err := b.Verify(); err != nil {
		t.Errorf("wanted no error, but got %v", err.Error())
	}
}

func assertTokenRequestVerifyError(t *testing.T, b tokenBuilder, wantedError string) {
	err := b.Verify()
	if err == nil {
		t.Errorf("wanted an error %v, but didn't get one", wantedError)
		return
	}

	if err.Error() != wantedError {
		t.Errorf("wanted an error %v, but got %v", wantedError, err.Error())
	}
}
