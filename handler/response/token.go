package response

// Presentation layer
// processes models for display
// convert models to json

import (
	"fmt"
	"net/http"
)

func ResponseRequestToken(w http.ResponseWriter, name string) {
	fmt.Fprint(w, name)
}
