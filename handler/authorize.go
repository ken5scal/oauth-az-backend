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
	return &authzHandler{r}
}

func (h *authzHandler) RequestAuthz(w http.ResponseWriter, r *http.Request) {

}

func (h *authzHandler) GenerateAuthzCode() (*domain.AuthorizationInfo, error) {
	r, _ := h.repo.(*infrastructure.AuthzInfoRepositoryImpl)
	r.BeginTransaction()

	azInfo := domain.AuthorizationInfoBuilder()
	err := h.repo.Insert(azInfo)

	if err != nil {
		r.Rollback()
		return nil, err
	}
	r.Commit()

	return azInfo, nil
}
