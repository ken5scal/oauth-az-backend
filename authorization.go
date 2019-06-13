package main

// Save this on DB, Redis/memchached
type AuthorizationInfo struct {
	AuthorizationId string // Not ID Token
	ClientId        string
	UserId          string
	Scope           []string
	RedirectUri     string
	AuthzCode       string // Expires in 10 min: https://tools.ietf.org/html/rfc6749#section-4.1.2
	RefreshToken    string
}

// Save this on Redis/memchached
type AccessToken struct {
	Token           string
	AuthorizationId string
	Expires         string
	IssuedDate      string
}

type UserInfo struct {
	UserId     string
	Name       string
	FamilyName string
	GivenName  string
}
