package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (s server) proxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, s.consolePrefix, http.StatusFound)
		return
	}

	token := s.proxyPathRegex.FindStringSubmatch(r.URL.Path)
	if len(token) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	upstreamId, found, err := s.dao.GetUpstreamIdByToken(token[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO
	_ = upstreamId
	upstream, _ := url.Parse("http://localhost:5555/")

	(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(upstream)
			r.Out.URL.Path = r.In.URL.Path[37:] // trim token prefix
			r.Out.Header.Del("Cookie")
		},
	}).ServeHTTP(w, r)
}
