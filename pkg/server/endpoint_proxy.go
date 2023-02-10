package server

import (
	"fmt"
	"net/http"
)

func (s server) proxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.redirectToConsole(w, r)
		return
	}

	fmt.Fprintln(w, "proxy()", r.URL.Path)
	// https://future-architect.github.io/articles/20230131a/
	// TODO: use httputil.ReverseProxy
}
