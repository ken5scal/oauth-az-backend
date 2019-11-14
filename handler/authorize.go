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
	azInfoBuilder := domain.AuthorizationInfoBuilder(nil)

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
