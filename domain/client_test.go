package domain

import (
	"testing"
)

func Test_client_RegisterRedirectUri(t *testing.T) {
	existingClient := ClientBuilder()
	existingClient.RedirectUris = []string{"http://localhost:8080"}
	tests := []struct {
		name    string
		client  *client
		args    []string
		wantErr bool
		errMsg  error
	}{
		{
			name:    "Single Registration",
			client:  ClientBuilder(),
			args:    []string{"http://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "Duplicated Registration to New Client",
			client:  ClientBuilder(),
			args:    []string{"http://localhost:8080", "http://localhost:8080"},
			wantErr: true,
			errMsg:  ErrDuplicatedRegistrationUris,
		},
		{
			name:    "Duplicated Registration to Existing Client",
			client:  existingClient,
			args:    []string{"http://localhost:8080"},
			wantErr: true,
			errMsg:  ErrDuplicatedRegistrationUris,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				ID:           tt.client.ID,
				Secrets:      tt.client.Secrets,
				RedirectUris: tt.client.RedirectUris,
				ClientType:   tt.client.ClientType,
				ClientStatus: tt.client.ClientStatus,
			}
			if err := c.RegisterRedirectUris(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("client.RegisterRedirectUri() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil && tt.wantErr) && (err.Error() != tt.errMsg.Error()) {
				t.Errorf("client.RegisterRedirectUri() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}
