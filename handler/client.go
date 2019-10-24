package handler

import (
	"fmt"
	"net/http"
)

func ClientServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "success")
}
