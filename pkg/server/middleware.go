package server

import (
	"fmt"
	"net/http"
	"net/url"
)

func method(method ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range method {
			if r.Method == m {
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func antiCSRF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref, err := url.Parse(r.Header.Get("Referer"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid Referer")
			return
		}
		if ref.Host == "" || r.Host != ref.Host {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "CSRF detected")
			return
		}
	}
}
