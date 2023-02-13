package handlers

import (
	"net/http"
	"net/http/httputil"
)

func (s server) proxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, s.consolePrefix, http.StatusFound)
		return
	}

	token := s.regex.FindStringSubmatch(r.URL.Path)
	if len(token) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	upstream, err := s.logic.GetUpstreamByToken(token[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if upstream == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: initialize at New()
	(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(upstream.Url)
			r.Out.URL.Path = r.In.URL.Path[37:] // trim token prefix
			r.Out.Header.Del("Cookie")
		},
	}).ServeHTTP(w, r)
}
