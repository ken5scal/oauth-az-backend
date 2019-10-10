package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"net/url"
)

var ErrDuplicatedRegistrationUris = errors.New("cannot register duplicated URIs")

type client struct {
	ID string
	// "minimum of 128 bits of entropy where the probability of an attacker guessing the generated token is
	// less than or equal to 2^(-160) as per [RFC6749] section 10.10"
	//  https://bitbucket.org/openid/fapi/pull-requests/45/bring-access-token-requirements-inline/diff
	// calculated by https://8gwifi.org/passwdgen.jsp
	Secrets      string
	RedirectUris []string
	ClientType   ClientType
	ClientStatus ClientStatus // manage

	// RP status
	//AuthzRevision int  // Let's not care at this point
}

func ClientBuilder() *client {
	var lengthEnoughForEntropy = 26
	return &client{
		ID:      uuid.New().String(),
		Secrets: generateSecret(lengthEnoughForEntropy),
		//AuthzRevision: 1,
	}
}

func generateSecret(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:length]
}

func (c *client) RegisterRedirectUris(uris []string) error {
	tmpUri := make(map[string]bool, len(uris))

	for _, uri := range uris {
		// Check All URIs formats are correct
		if _, err := url.Parse(uri); err != nil {
			return err
		}

		if _, exists := tmpUri[uri]; exists {
			return ErrDuplicatedRegistrationUris
		}

		tmpUri[uri] = false
	}

	for registeringUri, _ := range tmpUri {
		for _, existingUri := range c.RedirectUris {
			if registeringUri == existingUri {
				return ErrDuplicatedRegistrationUris
			}
		}
		c.RedirectUris = append(c.RedirectUris, registeringUri)
	}
	return nil
}

//func (c *client) CopyAuthzRevision() int {
//	return c.AuthzRevision
//}

// ClientType is the OAuth client types
// https://tools.ietf.org/html/rfc6749#section-2.1
type ClientType struct {
	value string
}

var Confidential = ClientType{"confidential"}
var Public = ClientType{"public"}

func (c ClientType) String() string {
	if c.value == "" {
		return "undefined client"
	}
	return c.value
}

// ClientStatus represents client's current availability
type ClientStatus struct {
	status string
}

var Developing = ClientStatus{"developing"}
var Published = ClientStatus{"published"}
var Suspended = ClientStatus{"suspended"}
var Deleted = ClientStatus{"deleted"}

func (c ClientStatus) String() string {
	if c.status == "" {
		return "developing"
	}
	return c.status
}
