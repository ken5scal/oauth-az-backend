package domain

import (
	"github.com/google/uuid"
)

type client struct {
	ID           string
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
