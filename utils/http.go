package utils

import (
	"fmt"
	"net/http"
)

func HttpError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, http.StatusText(code))
}
