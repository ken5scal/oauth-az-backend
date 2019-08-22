package handler

import (
	"github.com/ken5scal/oauth-az/domain"
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
