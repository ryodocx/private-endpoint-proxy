package server

import (
	"net/http"
	"net/http/httputil"
)

func (s server) proxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.redirectToConsole(w, r)
		return
	}

	token := s.regex.FindStringSubmatch(r.URL.Path)
	if len(token) < 2 || len(token[1]) != 36 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	upstream, err := s.dao.GetUpstreamByToken(token[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if upstream == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(upstream.Url)
			r.Out.URL.Path = r.In.URL.Path[37:] // trim token prefix
		},
	}).ServeHTTP(w, r)
}
