package main

type Client struct {
	ID         string
	Secret     string
	ClientType ClientType
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
