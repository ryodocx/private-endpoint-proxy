package server

import (
	"log"
	"net/http"
	"net/url"
)

func logging() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI, r.RemoteAddr, r.Header.Get("User-Agent"))
	}
}

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
