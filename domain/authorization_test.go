package domain

import (
	"testing"
	"time"
)

func TestAuthorizationInfo_isCodeValid(t *testing.T) {
	type fields struct {
		AuthorizationId string
		ClientId        string
		UserId          string
		Scope           []string
		RedirectUri     string
		AuthzCode       string
		CodeExpiration  time.Time
		RefreshToken    string
		AuthzRevision   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "whatever",
			fields: fields{
				CodeExpiration: time.Now().Add(time.Second * time.Duration(codeExpirationDuration)),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationInfo{
				AuthorizationId: tt.fields.AuthorizationId,
				ClientId:        tt.fields.ClientId,
				UserId:          tt.fields.UserId,
				Scope:           tt.fields.Scope,
				RedirectUri:     tt.fields.RedirectUri,
				AuthzCode:       tt.fields.AuthzCode,
				CodeExpiration:  tt.fields.CodeExpiration,
				RefreshToken:    tt.fields.RefreshToken,
				AuthzRevision:   tt.fields.AuthzRevision,
			}
			if got := a.isCodeValid(); got != tt.want {
				t.Errorf("AuthorizationInfo.isCodeValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
