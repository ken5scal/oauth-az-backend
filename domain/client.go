package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"net/url"
)

var ErrDuplicatedRegistrationUris = errors.New("cannot register duplicated URIs")
var ErrInvalidClientType = errors.New("cannot register client type other than confidential or public")

const (
	ConfidentialClient = "confidential"
	PublicClient       = "public"

	// "minimum of 128 bits of entropy where the probability of an attacker guessing the generated token is
	// less than or equal to 2^(-160) as per [RFC6749] section 10.10"
	//  https://bitbucket.org/openid/fapi/pull-requests/45/bring-access-token-requirements-inline/diff
	// calculated by https://8gwifi.org/passwdgen.jsp
	lengthEnoughForEntropy = 26
)

type client struct {
	id           string
	secrets      string
	redirectUris []string
	clientType   string

	// they are values not related to rfc, but friendly for managing authz endpoints
	name          string       // Not defined in rfc, but need for reaability
	clientStatus  ClientStatus // manage developing status
	authzRevision int
}

type builder struct {
	clientType string
}

func newClientBuilder() *builder {
	return &builder{}
}

func (cb *builder) ClientType(clientType clientType) *builder {
	cb.clientType = clientType.String()
	return cb
}

func (cb *builder) Build() (c *client, err error) {
	if !isClientTypeValid(cb.clientType) {
		return nil, ErrInvalidClientType
	}

	c = &client{
		id:         uuid.New().String(),
		clientType: cb.clientType,
		secrets:    generateSecret(lengthEnoughForEntropy),
	}

	return c, nil
}

func generateSecret(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:length]
}

// clientType is the OAuth client types
// https://tools.ietf.org/html/rfc6749#section-2.1
func isClientTypeValid(clientType string) bool {
	return clientType == ConfidentialClient || clientType == PublicClient
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
		for _, existingUri := range c.redirectUris {
			if registeringUri == existingUri {
				return ErrDuplicatedRegistrationUris
			}
		}
		c.redirectUris = append(c.redirectUris, registeringUri)
	}
	return nil
}

type clientType struct {
	clientType string
}

var confidential = clientType{"confidential"}
var public = clientType{"confidential"}

func (c clientType) String() string {
	return c.clientType
}

// clientStatus represents client's current availability
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
