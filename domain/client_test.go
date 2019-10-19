package domain

import (
	"testing"
	"testing/quick"
)

func Test_Client_Register(t *testing.T) {
	t.Run("register clietnt", func(t *testing.T) {
		got := NewClientBuilder().Build()
		if got.ID == "" || got.Secrets == "" {
			t.Errorf("got nothing, want %v", got)
		}
	})
}

func Test_client_RegisterType(t *testing.T) {
	builder := newClientBuilder()
	t.Run("register proper client type", func(t *testing.T) {
		properType := []string{"confidential", "public"}
		for _, clientType := range properType {
			_, err := builder.ClientType(clientType).Build()
			if err != nil {
				t.Errorf("got error for client type %v, want no error", clientType)
			}
		}
	})

	t.Run("register in-proper client type", func(t *testing.T) {
		got, err := builder.ClientType("not proper client type").Build()
		if got != nil {
			t.Error("wanted a nil client but didn't get one")
		}

		if err.Error() != ErrInvalidClientType.Error() {
			t.Errorf("got %q, want %q", got, ErrInvalidClientType.Error())
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

	//assertion := func(arabic uint16) bool {
	//	if arabic > 3999 {
	//		return true
	//	}
	//	log.Println(arabic)
	//	roman := ConvertToRoman(int(arabic))
	//	fromRoman := ConvertToArabic(roman)
	//	return uint16(fromRoman) == arabic
	//}
	//
	//quickConfig := &quick.Config{MaxCount:10}
	//
	//// quick.Check a function that it will run against a number of random inputs, if the function returns false it will be seen as failing the check.
	//if err := quick.Check(assertion, quickConfig); err != nil {
	//	t.Error("failed checks", err)
	//}
}

func Test_client_RegisterRedirectUri(t *testing.T) {
	existingClient := NewClientBuilder()
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
			client:  NewClientBuilder(),
			args:    []string{"http://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "Duplicated Registration to New Client",
			client:  NewClientBuilder(),
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
