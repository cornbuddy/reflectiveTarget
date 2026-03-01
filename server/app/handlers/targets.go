package handlers

import (
	"net/http"
)

type targetsHandler struct{}

func (h targetsHandler) list(w http.ResponseWriter, r *http.Request) {}

func (h targetsHandler) new(w http.ResponseWriter, r *http.Request) {}

func (h targetsHandler) update(w http.ResponseWriter, r *http.Request) {}
