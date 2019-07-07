package domain

import (
	"github.com/google/uuid"
)

type client struct {
	ID string
	// "minimum of 128 bits of entropy where the probability of an attacker guessing the generated token is less than or equal to 2^(-160) as per [RFC6749] section 10.10"
	//  https://bitbucket.org/openid/fapi/pull-requests/45/bring-access-token-requirements-inline/diff
	// calculated by https://8gwifi.org/passwdgen.jsp
	// Don't use symboles, just numbers and letters from 22 ~ 26
	Secrets      []string
	RedirectUris []string
	ClientType   ClientType
	ClientStatus ClientStatus // manage

	// RP status
	AuthzRevision int
}

func ClientBuilder() client {
	return client{
		ID:            uuid.New().String(),
		AuthzRevision: 1,
	}
}

func (c *client) GenerateSecret() {
	c.Secrets = append(c.Secrets, uuid.New().String())
}

func (c *client) CopyAuthzRevision() int {
	return c.AuthzRevision
}

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

type ClientStatus struct {
	status string
}

var Developing = ClientType{"developing"}
var Published = ClientType{"published"}
var Suspended = ClientType{"suspended"}
var Deleted = ClientType{"deleted"}

func (c ClientStatus) String() string {
	if c.status == "" {
		return "developing"
	}
	return c.status
}
