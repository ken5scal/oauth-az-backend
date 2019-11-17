package handler

import (
	"fmt"
	"github.com/ken5scal/oauth-az/domain"
	"net/http"
	"net/url"
	"reflect"
)

type authzHandler struct {
	repo domain.AuthzInfoRepository
}

func NewAuthzHandler(r domain.AuthzInfoRepository) *authzHandler {
	h := new(authzHandler)
	h.repo = r
	return h
}

// Check Duplicated query parameter
// https://tools.ietf.org/html/rfc6749#section-3.1.2
func isQueryParameterDuplicated(queryValues url.Values) bool {
	fmt.Println(queryValues)
	fmt.Println(len(queryValues.Get("redirect_uri")))
	fmt.Println(reflect.TypeOf(queryValues.Get("redirect_uri")).String())
	return reflect.TypeOf(queryValues.Get("redirect_uri")).String() == "[]string" ||
		reflect.TypeOf(queryValues.Get("client_id")).String() == "[]string" ||
		reflect.TypeOf(queryValues.Get("state")).String() == "[]string" ||
		reflect.TypeOf(queryValues.Get("response_type")).String() == "[]string"
}

// RequestAuthz assumes user-agent is a web browser
func (h *authzHandler) RequestAuthz(w http.ResponseWriter, r *http.Request) {
	// Check Authorization request
	params := r.URL.Query()
	clientId := params.Get("client_id")
	// Todo get client's registered redirection endpoint from data store
	azInfo, _ := h.repo.GetClientInfoByID(clientId)
	clientRedirectEps := []string{azInfo.RedirectUri}
	redirectUri, err := url.ParseRequestURI(params.Get("redirect_uri"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	builder := domain.AuthorizationInfoBuilder(params.Get("response_type"), clientId, params.Get("state"), redirectUri)
	if err := builder.Verify(clientRedirectEps); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	//if err != nil {
	// https://tools.ietf.org/html/rfc6749#section-4.1.2.1
	// ex:  HTTP/1.1 302 Found
	//   Location: https://client.example.com/cb?error=access_denied&state=xyz
	//w.WriteHeader(http.StatusFound)
	//w.Write([]byte(err.Error()))
	//}

	// Golang Automatically handle duplicated

	// TODO Do Authentication (Verify the identity of the resource owner) // https://tools.ietf.org/html/rfc6749#section-3.1
	// TODO Obtain Authorization Decision
	// TODO Directs the user-agent with Authorization Response
	// builder.Build
	//r, _ := h.repo.(*infrastructure.AuthzInfoRepositoryImpl)
	//r.BeginTransaction()
	//
	//azInfo := domain.AuthorizationInfoBuilder(nil)
	//err := h.repo.Insert(azInfo)
	//
	//if err != nil {
	//	r.Rollback()
	//	return nil, err
	//}
	//r.Commit()

	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	// ex: HTTP/1.1 302 Found
	//     Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA &state=xyz
	w.WriteHeader(http.StatusFound)
	w.Header().Add("Location", "https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz")
}
