package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

func (s server) console(w http.ResponseWriter, r *http.Request) {
	// strip path prefix
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.consolePrefix)

	// context
	u, err := s.upstream.Upstreams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := model.Context{
		User:      r.Header.Get("X-Forwarded-Email"),
		Upstreams: u,
		Query:     r.URL.Query(),
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
	tokens, err := s.dao.GetTokensByUserId(r.Header.Get("X-Forwarded-Email"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(tokens) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "you owned too many tokens")
		return
	}

	r.ParseForm()

	if err := s.dao.CreateToken(dao.Token{
		UserId:      r.Header.Get("X-Forwarded-Email"),
		Token:       uuid.NewString(),
		Description: r.PostFormValue("description"),
		UpstreamId:  r.PostFormValue("upstream"),
		CreatedAt:   time.Now(),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}

func (s server) deleteToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if err := s.dao.DeleteToken(r.PostFormValue("token")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal error")
		return
	}

	http.Redirect(w, r, s.consolePrefix, http.StatusFound)
}
