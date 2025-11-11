package handlers

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
)

type ShotsRouter struct {
	Validator validators.ShotsRequestValidator
}

func (r ShotsRouter) Post(resp http.ResponseWriter, req *http.Request) {
	http.Error(resp, "kek", http.StatusBadRequest)
}
