package domain

import (
	"net/url"
	"testing"
)

var wantedClientEndpoints = []string{"https://example.com/cb"}

func TestGeneratingInvalidAuthorizaionResponse(t *testing.T) {
	redirectUri, _ := url.ParseRequestURI(wantedClientEndpoints[0])

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	t.Run("request with empty responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("", "clientId", redirectUri)
		_, err := builder.Build(wantedClientEndpoints)
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	t.Run("request with unsupported responseType", func(t *testing.T) {
		builder := AuthorizationInfoBuilder("fake", "clientId", redirectUri)
		_, err := builder.Build(wantedClientEndpoints)
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", UnsupportedResponseType)
		} else if err.Error() != UnsupportedResponseType {
			t.Errorf("wanted an error %v, but got %v ", UnsupportedResponseType, err.Error())
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-3.1.2.3
	// https://tools.ietf.org/html/rfc6749#section-3.1.2.4
	t.Run("request with invalid redirectUri", func(t *testing.T) {
		builder := AuthorizationInfoBuilder(ResponseTypeCode, "clientId", redirectUri)
		badClientEndpoints := []string{"https://example.com/bad"}
		_, err := builder.Build(badClientEndpoints)
		if err == nil {
			t.Errorf("wanted an error %v, but didn't get one", InvalidRequest)
		} else if err.Error() != InvalidRequest {
			t.Errorf("wanted an error %v, but got %v ", InvalidRequest, err.Error())
		}
	})
}

func TestGeneratingAuthorizationCode(t *testing.T) {
	builder := &authorizationBuilder{
		responseType: ResponseTypeCode,
		clientId:     "clietndId",
	}
	requestedClientEP := wantedClientEndpoints[0]
	redirectUri, _ := url.ParseRequestURI(requestedClientEP)
	builder = AuthorizationInfoBuilder(ResponseTypeCode, "clientId", redirectUri)

	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	t.Run("generate simple auth code", func(t *testing.T) {
		az, err := builder.Build(wantedClientEndpoints)
		if err != nil {
			t.Errorf("got an error %v but didn't want one", err.Error())
		}

		if az.code == "" {
			t.Errorf("wanted a code but didn't get one")
		}
	})

	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	t.Run("generate auth code with state", func(t *testing.T) {
		az, _ := builder.State("xyz").Build(wantedClientEndpoints)
		if az.state == "" && az.state != builder.state {
			t.Errorf("wanted a state %v, but got %v", az.state, builder.state)
		}
	})

	// // https://tools.ietf.org/html/rfc6749#section-4.1.2
	t.Run("generate auth code with redirect uri", func(t *testing.T) {
		az, _ := builder.Build(wantedClientEndpoints)
		if az.redirectUri.Scheme+"://"+az.redirectUri.Host+az.redirectUri.Path != requestedClientEP {
			if az.redirectUri.RawQuery == "" {
				t.Errorf("wanted a redirect Uri, but didn't get one")
			}
			t.Errorf("got redirect Uri %v but wanted %v", az.redirectUri, redirectUri)
		}
	})
}
