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
		fail := func() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("CSRF detected"))
		}

		ref, err := url.Parse(r.Header.Get("Referer"))
		if err != nil {
			fail()
			return
		}

		if r.Host != ref.Host {
			fail()
			return
		}
	}
}
