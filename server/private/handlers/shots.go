package handlers

import (
	"net/http"
)

type ShotsRouter struct{}

func (r ShotsRouter) Post(resp http.ResponseWriter, req *http.Request) {
	http.Error(resp, "kek", http.StatusBadRequest)
}
