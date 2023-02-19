package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

func (s server) console(w http.ResponseWriter, r *http.Request) {
	// strip path prefix
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.consolePrefix)

	// context
	upstreams, err := s.config.Upstreams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := model.Context{
		User:      r.Header.Get("X-Forwarded-Email"),
		Upstreams: upstreams,
	}
	tokens, err := s.dao.GetTokensByUserId(ctx.User)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i, t := range tokens {
		u, _ := url.Parse(fmt.Sprintf("https://example.com/%d", i))

		upstream := model.Upstream{
			Id:          fmt.Sprintf("upstream%d", i),
			Description: fmt.Sprintf("Description%d", i),
			Url:         u,
		}
		ctx.Tokens = append(ctx.Tokens, t.ToModel(upstream))
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
	return
}

func (s server) addToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	log.Println(r.PostForm.Encode())

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}

func (s server) deleteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	log.Println(r.PostForm.Encode())

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}
