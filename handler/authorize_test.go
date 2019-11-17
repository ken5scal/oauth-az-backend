package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var responseType = "code"
var clientId = "s6BhdRkqt3"
var state = "xyz"
var redirectUri = "https://client.example.com/cb"

func TestHandlingInvalidAuthorizationRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	response := httptest.NewRecorder()

	q := url.Values{}
	//q.Add("response_type", "fake"+responseType)
	q.Add("client_id", clientId)
	q.Add("state", state)
	q.Add("redirect_uri", redirectUri)

	server := http.NewServeMux()
	server.HandleFunc("/authorize", NewAuthzHandler(&DummyRepository{}).RequestAuthz)

	t.Run("request with empty response type", func(t *testing.T) {
		request.URL.RawQuery = q.Encode()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Errorf("wanted http statsu code %v, but got %v", http.StatusBadRequest, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

	t.Run("request with unsupported response type", func(t *testing.T) {
		q.Add("response_type", "fake"+responseType)
		request.URL.RawQuery = q.Encode()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Errorf("wanted http statsu code %v, but got %v", http.StatusBadRequest, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

	t.Run("request with broken redirect uri", func(t *testing.T) {
		q.Set("response_type", responseType)
		q.Set("redirect_uri", "i'm broken redirect uri")
		request.URL.RawQuery = q.Encode()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusInternalServerError {
			t.Errorf("wanted http statsu code %v, but got %v", http.StatusInternalServerError, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

}

func TestAuthorizationHeader(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	response := httptest.NewRecorder()

	t.Run("returns Authorization error response if request is invalid", func(t *testing.T) {

	})

	t.Run("returns Authorization success response", func(t *testing.T) {
		requiredParams := map[string]string{
			"response_type": "code",
			"client_id":     "s6BhdRkqt3",
			"redirect_uri":  "https://client.example.com/cb",
			//scope
			//state
		}
		t.Errorf("did not get correct status, got %d, want %d", response.Code, http.StatusFound)
		q := request.URL.Query()
		for k, v := range requiredParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()

		server := http.NewServeMux()
		handler := NewAuthzHandler(&DummyRepository{})
		server.HandleFunc("/authorize", handler.RequestAuthz)
		server.ServeHTTP(response, request)

		u, err := url.Parse(response.Header().Get("Location"))
		if err != nil {
			t.Errorf("got an error when parsing authorization response's location header: %v", err)
		}

		redirectUri := u.Scheme + "://" + u.Host + u.Path
		if redirectUri != requiredParams["redirect_uri"] {
			t.Errorf("did not get correct redirect url, got %v, want %v", redirectUri, requiredParams["redirect_uri"])
		}

		if u.Query().Get("code") == "" {
			t.Error("code parameter in authorization response is required")
		}

		// state is required iff state parameter was present in the request
		// RFC 6749 RECOMMENDS it
		if q.Get("state") != "" {
			state := u.Query().Get("state")
			if state == "" {
				t.Error("state parameter in authorization response is required")
			}

			if state != q.Get("state") {
				t.Errorf("did not get correct status, got %s, want %s", state, q.Get("state"))
			}
		}
	})
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
