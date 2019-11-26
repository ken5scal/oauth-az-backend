package domain

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

var registeredClientEndpoints = []string{"https://example.com/cb"}
var state = "xyz"

func assertError(t *testing.T, got, want *authorizationError) {
	t.Helper()
	if want.state != "" {
		if got.state == "" {
			t.Errorf("wanted a state %v, but didn't get one", want.state)
		} else if got.state != want.state {
			t.Errorf("wanted a state %v, but got %v ", want.state, got.state)
		}
	}

	if got == nil {
		t.Errorf("wanted an error %v, but didn't get one", want.error)
	} else if got.Error() != want.Error() {
		t.Errorf("wanted an error %v, but got %v ", want.error, got.Error())
	}
}

func TestInvalidAuthorizationRequest(t *testing.T) {
	validRedirectUri, _ := url.ParseRequestURI(registeredClientEndpoints[0])

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("", "clientId", state, validRedirectUri)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), state})
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	t.Run("request with unsupported responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("fake", "clientId", state, validRedirectUri)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(UnsupportedResponseType), state})
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty clientId", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "", "xyz", validRedirectUri)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), state})
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty state", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "", validRedirectUri)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), ""})
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.1
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.2
	t.Run("request with empty redirectUri", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", nil)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), state})
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2
	t.Run("request with redirect uri contains fragment component", func(t *testing.T) {
		badRedirectUriParam, _ := url.ParseRequestURI(registeredClientEndpoints[0] + "#fragment")
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", badRedirectUriParam)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), state})
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2.3
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.4
	t.Run("request with unregistered redirect uri", func(t *testing.T) {
		badRedirectUriParam, _ := url.ParseRequestURI("https://example.com/bad")
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", badRedirectUriParam)
		err := builder.Verify(registeredClientEndpoints)
		assertError(t, err, &authorizationError{errors.New(InvalidRequest), state})
	})
}

func TestValidAuthorizationRequest(t *testing.T) {
	requestedClientEP := registeredClientEndpoints[0]
	redirectUri, _ := url.ParseRequestURI(requestedClientEP)
	builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", redirectUri)

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	t.Run("check simple auth request", func(t *testing.T) {
		if err := builder.Verify(registeredClientEndpoints); err != nil {
			t.Errorf("got an error %v but didn't want one", err.Error())
		}
		if builder.Build().code == "" {
			t.Errorf("wanted a code but didn't get one")
		}

		builder.responseType = responseTypeCodeIdToken
		if err := builder.Verify(registeredClientEndpoints); err != nil {
			t.Errorf("got an error %v but didn't want one", err.Error())
		}

		az := builder.Build()
		localTime := time.Now().Local()

		if az.code == "" {
			t.Errorf("wanted a code but didn't get one")
		}

		if az.codeExpiration.Sub(localTime) > time.Minute*time.Duration(codeExpirationDuration) {
			t.Errorf("wanted authorization code lifetime to be less than 10 minutes, but got %v", az.codeExpiration.Sub(localTime))
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	t.Run("check state holds the same value", func(t *testing.T) {
		az := builder.Build()
		if az.state == "" && az.state != builder.state {
			t.Errorf("wanted a state %v, but got %v", az.state, builder.state)
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2.2
	t.Run("check redirect uri is authz request is pre-registered", func(t *testing.T) {
		az := builder.Build()
		if az.redirectUri.Scheme+"://"+az.redirectUri.Host+az.redirectUri.Path != requestedClientEP {
			if az.redirectUri.RawQuery == "" {
				t.Errorf("wanted a redirect Uri, but didn't get one")
			}
			t.Errorf("got redirect Uri %v but wanted %v", az.redirectUri, redirectUri)
		}
	})
}

func Test_authorization_ReturnRedirectionEndpoint(t *testing.T) {
	u, _ := url.ParseRequestURI("https://client.example.com/cb")
	az := &authorization{
		code:        "SplxlOBeZQQYbYS6WxSbIA",
		state:       "xyz",
		redirectUri: u,
	}

	got := az.ReturnRedirectionEndpoint()
	expected := "https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz"

	if got != expected {
		t.Errorf("ReturnRedirectionEndpoint() = %v, want %v", got, expected)
	}
}
