package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"net/http"
	"net/url"
)

const (
	authzRequestParamRedirectUri  = "redirect_uri"
	authzRequestParamCode         = "code"
	authzRequestParamState        = "state"
	authzRequestParamResponseType = "response_type"
	authzRequestParamClientId     = "client_id"
)

type authzHandler struct {
	repo domain.AuthzInfoRepository
}

func NewAuthzHandler(r domain.AuthzInfoRepository) *authzHandler {
	h := new(authzHandler)
	h.repo = r
	return h
}

// RequestAuthz assumes user-agent is a web browser
func (h *authzHandler) RequestAuthz(w http.ResponseWriter, r *http.Request) {
	// Check Authorization request
	params := r.URL.Query()

	// Todo get client's registered redirection endpoint from data store
	clientId := params.Get(authzRequestParamClientId)
	azInfo, _ := h.repo.GetClientInfoByID(clientId)
	clientRedirectEps := []string{azInfo.RedirectUri}

	redirectUri, err := url.ParseRequestURI(params.Get(authzRequestParamRedirectUri))

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
	w.WriteHeader(http.StatusFound)
	w.Header().Add("Location", builder.Build().ReturnRedirectionEndpoint())
}
