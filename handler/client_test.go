package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//type clientHandler struct {
//	repo domain.ClientRepository
//}

// registerClient
// unRegisterClient
// suspendClient
// releaseClient

func TestGETClients(t *testing.T) {
	t.Run("registers client", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/client/org/", nil)
		response := httptest.NewRecorder()

		ClientServer(response, request)

		got := response.Body.String()
		want := "client-a"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
