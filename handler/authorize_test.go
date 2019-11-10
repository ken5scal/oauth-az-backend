package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizationHeader(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	response := httptest.NewRecorder()

	requiredParams := map[string]string{
		"response_type": "code",
		"client_id":     "s6BhdRkqt3",
		"redirect_uri":  "https://client.example.com/cb",
		//scope
		//state
	}

	q := request.URL.Query()
	for k, v := range requiredParams {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()

	server := http.NewServeMux()
	handler := NewAuthzHandler(&DummyRepository{})
	server.HandleFunc("/authorize", handler.RequestAuthz)
	server.ServeHTTP(response, request)

	if response.Code != http.StatusFound {
		t.Errorf("did not get correct status, got %d, want %d", response.Code, http.StatusFound)
	}
}

type DummyRepository struct {
}

func (r *DummyRepository) GetAuthzInfoForAccessToken(clientID, userID string) (*domain.AuthorizationInfo, error) {
	return nil, nil
}

func (r *DummyRepository) GetAuthzInfoByID(authzInfoID string) (*domain.AuthorizationInfo, error) {
	return nil, nil
}

func (r *DummyRepository) Insert(t *domain.AuthorizationInfo) error {
	return nil
}

func (r *DummyRepository) Update(t *domain.AuthorizationInfo) error {
	return nil
}

func (r *DummyRepository) Delete(t *domain.AuthorizationInfo) error {
	return nil
}
