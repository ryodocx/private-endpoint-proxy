package handlers

import (
	"fmt"
	"net/http"
)

func (s server) live(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "live")
}

func (s server) ready(w http.ResponseWriter, r *http.Request) {
	// TODO
	fmt.Fprintln(w, "ready")
}
