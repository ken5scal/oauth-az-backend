package main

// AuthorizationInfo is authorized info that RO granted to the RP
// so this is tied to user
// Create, Read, Update, Delete
// Save this on Redis/memchached, but also in long term DB
// because every authCode in token request are required to
// be compared with past one
// This can be deleted once user revoked the app
type AuthorizationInfo struct {
	AuthorizationId string // Not ID Token, should be indexed in DB
	ClientId        string // should be indexed combined with UserId
	UserId          string
	Scope           []string
	RedirectUri     string
	// Expires in 10 min: https://tools.ietf.org/html/rfc6749#section-4.1.2
	AuthzCode     string // should be indexed in DB
	RefreshToken  string // should be indexed in DB
	AuthzRevision int
}

func AuthorizationInfoBuilder(c *client) AuthorizationInfo {
	return AuthorizationInfo{
		AuthzRevision: c.AuthzRevision,
	}
}

// If Authorization Info does not exist, then the app sends an authorization request for the first time.
// If scope is changed and more values are specified, needs to ask for new consent
// If scope is decreased or unched, then't consent is not required
func (a *AuthorizationInfo) isConsentNeeded() bool {
	return false
	// if this is true, then AS needs considerring for
	// revoking access token, or refresh tokens
}

func (a *AuthorizationInfo) isRefreshTokenValid() bool {
	return false
	// check validity of Authorization Info
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

func (a *AccessTokenInfo) isAccessTokenValid() {
	// check token info              -> 401 invalid_token, "the access token is invalid"
	// check revocation status       -> 401 invalid_token, "the access token is invalid"
	// check authorization info      -> 401 invalid_token, "the access token is invalid"
	// check client status           -> 401 invalid_client
	// check user status             -> 401 invalid_user
	// check scope                   -> 401 insufficient_scope
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

// Revoke Access Token based om user request
func (a *AccessTokenInfo) userRevoke() {
	// 1. revoke access token (ofcourse)
	// 2. revoke refresh token (of course), any token request using refresh token should be also be invalid (so don't erase)
	// 3. also DELETE authorization info
}

// Revoke All Access Token tied to the RP by AS admins
func (a *AccessTokenInfo) adminRevoke() {
	// 1. revoke all access token issued for the RP
	// 2. revoke all refresh token issued for the RP
	// 3. set Client Status to be suspended

	// Aside from this function...
	// 4. if things are solved, update revision in Client
	// 5. change buck the client status to "published"
}

// if this is false, then Authorization Info is old
func (a *AccessTokenInfo) isAuthzRevisionValid() bool {
	// get Authorization Info from a.AuthorizationId
	// get Client Info from authZInfo.ClientId
	// compare authorizationInfo.authz_revision and client.authz_revision
	// if false, then token should be invalid

	return false
}

func (a *AccessTokenInfo) verifyToken() bool {
	return false //isRevoked  // 401 Unauthorized
}

func (a *AccessTokenInfo) GetTokenInfo() {
	// return a.AuthorizationId.ClientId, a.ExpiresIn
}

type UserInfo struct {
	UserId     string
	Name       string
	FamilyName string
	GivenName  string
}

func issueAccessToken() {

}
