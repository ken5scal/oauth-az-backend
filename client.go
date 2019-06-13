package main

import (
	"github.com/google/uuid"
)

type client struct {
	ID           string
	Secrets      []string
	RedirectUris []string
	ClientType   ClientType
	ClientStatus string
}

func ClientBuilder() client {
	return client{
		ID: uuid.New().String(),
	}
}

func (c *client) GenerateSecret() {
	c.Secrets = append(c.Secrets, uuid.New().String())
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
