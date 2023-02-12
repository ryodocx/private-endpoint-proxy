package server

import (
	"log"
	"net/http"
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
			w.Write([]byte("invalid Referer"))
		}

		referer := r.Header.Get("Referer")
		log.Println("referer: ", referer)
		if referer == "" {
			fail()
			return
		}

		// TODO: check referrer more strict
	}
}
