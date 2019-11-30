package handler

import (
	"fmt"
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
var returnFakeUri = false

func TestHandlingInvalidAuthorizationRequest(t *testing.T) {
	returnFakeUri = true
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	response := httptest.NewRecorder()

	q := url.Values{}
	q.Add(authzRequestParamClientId, clientId)
	q.Add(authzRequestParamState, state)
	q.Add(authzRequestParamRedirectUri, redirectUri)
	request.URL.RawQuery = q.Encode()

	server := http.NewServeMux()
	server.HandleFunc("/authorize", NewAuthzHandler(&DummyRepository{}).RequestAuthz)
	server.ServeHTTP(response, request)

	t.Run("request with empty response type", func(t *testing.T) {
		if response.Code != http.StatusBadRequest {
			t.Errorf("wanted http statsu code %v, but got %v", http.StatusBadRequest, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

	t.Run("request with unsupported response type", func(t *testing.T) {
		q.Add(authzRequestParamResponseType, "fake"+responseType)
		request.URL.RawQuery = q.Encode()
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Errorf("wanted http status code %v, but got %v", http.StatusBadRequest, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

	t.Run("request with invalid redirect uri", func(t *testing.T) {
		q.Set(authzRequestParamResponseType, responseType)
		q.Set(authzRequestParamRedirectUri, "i'm broken redirect uri")
		request.URL.RawQuery = q.Encode()
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)
		if response.Code != http.StatusInternalServerError {
			t.Errorf("wanted http status code %v, but got %v", http.StatusInternalServerError, response.Code)
		}
		if response.Body == nil {
			t.Error("wanted an error in response, but got none")
		}
	})

	t.Run("request with duplicated parameters", func(t *testing.T) {
		request.URL.RawQuery = request.URL.RawQuery + "&" + authzRequestParamRedirectUri + "=fakeuri"
		fmt.Println(request.URL.RawQuery)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		fmt.Println(response.Code)
		if response.Code != http.StatusBadRequest {
			t.Errorf("wanted http status code %v, but got %v", http.StatusBadRequest, response.Code)
		}
	})
}

func TestAuthorizationHeader(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/authorize", nil)
	response := httptest.NewRecorder()

	q := url.Values{}
	q.Add(authzRequestParamClientId, clientId)
	q.Add(authzRequestParamState, state)
	q.Add(authzRequestParamRedirectUri, redirectUri)
	q.Add(authzRequestParamResponseType, "code")
	request.URL.RawQuery = q.Encode()

	server := http.NewServeMux()
	server.HandleFunc("/authorize", NewAuthzHandler(&DummyRepository{}).RequestAuthz)

	t.Run("returns Authorization success response", func(t *testing.T) {
		server.ServeHTTP(response, request)

		if response.Code != http.StatusFound {
			t.Errorf("got http statsu code %v, but wanted %v", response.Code, http.StatusFound)
		}

		u, err := url.Parse(response.Header().Get("Location"))
		if err != nil {
			t.Errorf("got an error when parsing authorization response's location header: %v", err)
		}

		redirectUriInResponse := u.Scheme + "://" + u.Host + u.Path
		if redirectUriInResponse != redirectUri {
			t.Errorf("got redirect uri %v, but wanted %v", redirectUriInResponse, redirectUri)
		}

		// https://tools.ietf.org/html/rfc6749#section-4.1.2
		if u.Query().Get(authzRequestParamCode) == "" {
			t.Error("code parameter in authorization response is required")
		}

		if u.Query().Get("error") != "" {
			t.Error("got an error parameter but didn't want one")
		}

		// https://tools.ietf.org/html/rfc6749#section-4.1.2
		if q.Get(authzRequestParamState) != "" {
			if state := u.Query().Get(authzRequestParamState); state == "" {
				t.Error("state parameter in authorization response is required")
			} else if state != q.Get(authzRequestParamState) {
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

func (r *DummyRepository) GetClientInfoByID(authzInfoID string) (*domain.AuthorizationInfo, error) {
	u, _ := url.Parse(redirectUri)
	fu, _ := url.Parse("fake" + redirectUri)
	azInfo := &domain.AuthorizationInfo{RedirectUri: u}
	if returnFakeUri {
		azInfo.RedirectUri = fu
		return azInfo, nil
	}
	return azInfo, nil
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
