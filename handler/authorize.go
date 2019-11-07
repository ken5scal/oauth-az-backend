package handler

import (
	"github.com/ken5scal/oauth-az/domain"
	"github.com/ken5scal/oauth-az/infrastructure"
	"net/http"
)

const authorizationRequestMediaType = "application/x-www-form-urlencoded"

type authzHandler struct {
	repo domain.AuthzInfoRepository
}

func NewAuthzHandler(r domain.AuthzInfoRepository) *authzHandler {
	h := new(authzHandler)
	h.repo = r

	return h
}

func (h *authzHandler) RequestAuthz(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != authorizationRequestMediaType {

	}
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
