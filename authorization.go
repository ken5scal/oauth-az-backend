package main

// AuthorizationInfo is authorized info that RO granted to the RP
// Create, Read, Update, Delete
// Save this on Redis/memchached, but also in long term DB
// because every authCode in token request are required to
// be compared with past one
type AuthorizationInfo struct {
	AuthorizationId string // Not ID Token, should be indexed in DB
	ClientId        string // should be indexed combined with UserId
	UserId          string
	Scope           []string
	RedirectUri     string
	// Expires in 10 min: https://tools.ietf.org/html/rfc6749#section-4.1.2
	AuthzCode    string // should be indexed in DB
	RefreshToken string // should be indexed in DB
}

// If Authorization Info does not exist, then the app sends an authorization request for the first time.
// If scope is changed and more values are specified, needs to ask for new consent
// If scope is decreased or unched, then't consent is not required
func (a *AuthorizationInfo) isConsentNeeded() bool {
	return false
	// if this is true, then AS needs considerring for
	// revoking access token, or refresh tokens
}

// AccessTokenInfo is token issued based on AuthorizationInfo
// Create, Read, Update, Delete
// Save this on Redis/memchached
// If this is revoked, should eliminate them
// so dont's save it in long term db
type AccessTokenInfo struct {
	Token           string // AccessTokenInfo must be searchable from AccessToken
	AuthorizationId string // AccessTokenInfo must be searchable from AuthorizationId (like when user cancelled delegation)
	ExpiresIn       string // in seconds
	IssuedDate      string
}

func (a *AccessTokenInfo) isExpired() bool {
	return false
	// if true, then remove acccess token from storage
	// or just set expiration timing in Redis/Memchaced the same date as the expires_in
	// or just implement expires_in info in access token it self and compares
}

func (a *AccessTokenInfo) revoke() {
	// Check Expiration

	// User cancelled authorization

	// Revoked by us

	// Revoked by developers
}

// Revoke Access Token automatically bbased on `expires_in`
func (a *AccessTokenInfo) autoRevoke() {
	// check expiration and revoke the access token.
	// calculate current time from issued date and expire in
	// received token from RP will be compared to saved access token
	// if the token is expired , return an error

	// using JWT, access token will `self-contain` the expiration date

	// Save info in memcached/redis
	// expiration can be set by

	// also must erase access token
}

func (a *AccessTokenInfo) verifyToken() bool {
	return false //isRevoked  // 401 Unauthorized
}

type UserInfo struct {
	UserId     string
	Name       string
	FamilyName string
	GivenName  string
}

func issueAccessToken() {

}
