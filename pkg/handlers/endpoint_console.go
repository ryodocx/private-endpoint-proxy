package handlers

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/ryodocx/private-endpoint-proxy/pkg/interfaces"
)

func (s server) redirectToConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.consolePrefix, http.StatusTemporaryRedirect)
}

type context struct {
	User      string
	Tokens    []*interfaces.Token
	Upstreams []*interfaces.Upstream
}

func (s server) console(w http.ResponseWriter, r *http.Request) {
	// strip path prefix
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.consolePrefix)

	ctx := context{
		User: r.Header.Get("X-Forwarded-Email"),
	}

	if v, err := s.dao.GetTokens(ctx.User); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		ctx.Tokens = v
	}
	if v, err := s.dao.GetUpstreams(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		ctx.Upstreams = v
	}

	render := func(path string) {
		t, ok := s.templates[path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		var buf bytes.Buffer
		if err := t.Execute(&buf, ctx); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
		return
	}

	// index.html
	if r.URL.Path == "" {
		render("index.html")
		return
	}

	// static files
	if strings.Contains(r.URL.Path, ".") {
		http.FileServer(s.files).ServeHTTP(w, r)
		return
	}

	// templates
	render(r.URL.Path + ".html")
}

func (s server) addToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: anti CSRF

	r.ParseForm()
	log.Println(r.PostForm.Encode())

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}

func (s server) deleteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: anti CSRF

	r.ParseForm()
	log.Println(r.PostForm.Encode())

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}
