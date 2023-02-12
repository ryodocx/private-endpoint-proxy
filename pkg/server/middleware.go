package server

import (
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
			w.Write([]byte("invalid Referer"))
			return
		}
		if ref.Host == "" || r.Host != ref.Host {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("CSRF detected"))
			return
		}
	}
}
