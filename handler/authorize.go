package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"github.com/ken5scal/oauth-az/infrastructure"
	"net/http"
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
	az, err := domain.AuthorizationInfoBuilder(nil).Build(clietntRedirectEPs)
	if err != nil {
		// https://tools.ietf.org/html/rfc6749#section-4.1.2.1
		// ex:  HTTP/1.1 302 Found
		//   Location: https://client.example.com/cb?error=access_denied&state=xyz
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(err.Error()))
	}

	// TODO Do Authentication (Verify the identity of the resource owner) // https://tools.ietf.org/html/rfc6749#section-3.1
	// TODO Obtain Authorization Decision
	// TODO Directs the user-agent with Authorization Response

	// https://tools.ietf.org/html/rfc6749#section-4.1.2
	// ex: HTTP/1.1 302 Found
	//     Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA &state=xyz
	w.WriteHeader(http.StatusFound)
	w.Header().Add("Location", "https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz")
}

func (h *authzHandler) GenerateAuthzCode() (*domain.AuthorizationInfo, error) {
	r, _ := h.repo.(*infrastructure.AuthzInfoRepositoryImpl)
	r.BeginTransaction()

	azInfo := domain.AuthorizationInfoBuilder(nil)
	err := h.repo.Insert(azInfo)

	if err != nil {
		r.Rollback()
		return nil, err
	}
	r.Commit()

	return azInfo, nil
}
