package domain

import (
	"net/url"
	"testing"
	"time"
)

var registeredClientEndpoints = []string{"https://example.com/cb"}

func TestInvalidAuthorizationRequest(t *testing.T) {
	validRedirectUri, _ := url.ParseRequestURI(registeredClientEndpoints[0])

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("", "clientId", "xyz", validRedirectUri)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	t.Run("request with unsupported responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("fake", "clientId", "xyz", validRedirectUri)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", UnsupportedResponseType)
		} else if err.Error() != UnsupportedResponseType {
			t.Errorf("wanted an error %v, but got %v ", UnsupportedResponseType, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty clientId", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "", "xyz", validRedirectUri)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	t.Run("request with empty state", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "", validRedirectUri)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.1
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.2
	t.Run("request with empty redirectUri", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", nil)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2
	t.Run("request with redirect uri contains fragment component", func(t *testing.T) {
		badRedirectUriParam, _ := url.ParseRequestURI(registeredClientEndpoints[0] + "#fragment")
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", badRedirectUriParam)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2.3
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.4
	t.Run("request with unregistered redirect uri", func(t *testing.T) {
		badRedirectUriParam, _ := url.ParseRequestURI("https://example.com/bad")
		builder := AuthorizationInfoBuilder(responseTypeCode, "clientId", "xyz", badRedirectUriParam)
		if err := builder.Verify(registeredClientEndpoints); err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
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
		az := builder.State("xyz").Build()
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
