package handlers

import (
	"net/http"
	"net/http/httputil"
	"path"
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

	upstream, ok, err := s.upstream.Upstream(upstreamId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(upstream.Url)
			r.Out.URL.Path = path.Join(upstream.Url.Path, r.In.URL.Path[37:]) // trim token prefix
			r.Out.Header.Del("Cookie")
		},
	}).ServeHTTP(w, r)
}
