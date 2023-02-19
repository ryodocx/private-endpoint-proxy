package handlers

import (
	"fmt"
	"net/http"
)

func (s server) live(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "live")
}

func (s server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.upstream.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "upstream provider not ready")
		return
	}
	if err := s.dao.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "dao not ready")
		return
	}

	fmt.Fprintln(w, "ready")
}
