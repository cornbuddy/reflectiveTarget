package utils

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

var (
	htmxRequestKey = http.CanonicalHeaderKey("HX-Request")
)

// distinguish is it htmx request, and if it is, passes htmx specific headers to
// the response. othervise, performs standard http redirect. see
// https://stackoverflow.com/a/77252592
func Redirect(w http.ResponseWriter, r *http.Request, url, body string) {
	isHtmx := r.Header.Get(htmxRequestKey) == "true"
	if isHtmx {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

	log := log.Logger(r.Context())
	LogBadWrites(log)(w.Write([]byte(body)))
}
