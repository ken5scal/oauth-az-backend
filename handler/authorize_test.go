package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizationHeader(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	request.Header.Set("Content-Type", "fake"+authorizationRequestMediaType)

	router := http.NewServeMux()
	router.HandleFunc("/authorize", NewAuthzHandler(&DummyRepository{}).RequestAuthz)
	response := httptest.NewRecorder()

	if response.Code != http.StatusUnsupportedMediaType {
		t.Errorf("got response code %d want %d", response.Code, http.StatusUnsupportedMediaType)
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
