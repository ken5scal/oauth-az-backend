package domain

import (
	"testing"
	"testing/quick"
)

func Test_client_Register_Name(t *testing.T) {
	builder := newClientBuilder()
	t.Run("register client without name", func(t *testing.T) {
		_, err := builder.ClientType(confidential).Build()
		if err != nil {
			t.Error("wanted a nil error but got one")
		}
	})
}

func Test_client_Register_ClientType(t *testing.T) {
	builder := newClientBuilder()
	t.Run("register proper client type", func(t *testing.T) {
		properType := []clientType{confidential, public}
		for _, clientType := range properType {
			_, err := builder.ClientType(clientType).Build()
			if err != nil {
				t.Errorf("got error for client type %v, want no error", clientType)
			}
		}
	})

	t.Run("register in-proper client type", func(t *testing.T) {
		got, err := builder.ClientType(clientType{}).Build()
		if got != nil {
			t.Error("wanted a nil but got a client")
		}

		if err.Error() != ErrInvalidClientType.Error() {
			t.Errorf("got %v, want %v", got, ErrInvalidClientType.Error())
		}
	})

	t.Run("register in-proper client type", func(t *testing.T) {
		assertion := func(clientType string) bool {
			return !isClientTypeValid(clientType)
		}
		quickConfig := &quick.Config{MaxCount: 10}
		if err := quick.Check(assertion, quickConfig); err != nil {
			t.Error("failed checks", ErrInvalidClientType)
		}
	})
}

func Test_client_RegisterRedirectUri(t *testing.T) {
	existingClient, _ := newClientBuilder().ClientType(confidential).Build()
	testClient1, _ := newClientBuilder().ClientType(confidential).Build()
	testClient2, _ := newClientBuilder().ClientType(confidential).Build()
	existingClient.redirectUris = []string{"http://localhost:8080"}
	tests := []struct {
		name    string
		client  *client
		args    []string
		wantErr bool
		errMsg  error
	}{
		{
			name:    "Single Registration",
			client:  testClient1,
			args:    []string{"http://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "Duplicated Registration to New Client",
			client:  testClient2,
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
				id:           tt.client.id,
				secrets:      tt.client.secrets,
				redirectUris: tt.client.redirectUris,
				clientType:   tt.client.clientType,
				clientStatus: tt.client.clientStatus,
			}
			if err := c.RegisterRedirectUris(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("client.RegisterRedirectUri() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil && tt.wantErr) && (err.Error() != tt.errMsg.Error()) {
				t.Errorf("client.RegisterRedirectUri() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}
