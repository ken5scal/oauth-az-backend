package domain

import (
	"net/url"
	"testing"
)

func TestGeneratingAuthorizationCode(t *testing.T) {
	builder := &authorizationBuilder{
		responseType: "code",
		clientId:     "clietndId",
	}

	t.Run("generate simple auth code", func(t *testing.T) {
		az, err := builder.Build()
		if err != nil {
			t.Errorf("got an error %v but didn't want one", err.Error())
		}

		if az.code == "" {
			t.Errorf("wanted a code but didn't get one")
		}
	})

	t.Run("generate auth code with redirect uri", func(t *testing.T) {
		redirectUri, _ := url.Parse("https://client.example.com/cb")
		az, _ := builder.RedirectUri(redirectUri).Build()
		if az.redirectUri != redirectUri {
			if az.redirectUri.RawQuery == "" {
				t.Errorf("wanted a redirect Uri, but didn't get one")
			}
			t.Errorf("got redirect Uri %v but wanted %v", az.redirectUri, redirectUri)
		}
	})

	t.Run("generate auth code with state", func(t *testing.T) {
		az, _ := builder.State("xyz").Build()
		if az.state == "" && az.state != builder.state {
			t.Errorf("wanted a state %v, but got %v", az.state, builder.state)
		}
	})
}
