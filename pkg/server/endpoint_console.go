package server

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s server) redirectToConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.consolePrefix, http.StatusTemporaryRedirect)
}

type context struct {
	User      string
	Tokens    []token
	Upstreams []upstream
}
type token struct {
	Token       string
	Description string
	Upstream    string
	CreatedAt   time.Time
}
type upstream struct {
	Description string
	Id          string
}

func (s server) console(w http.ResponseWriter, r *http.Request) {
	// strip path prefix
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.consolePrefix)

	data := context{
		User: r.Header.Get("X-Forwarded-Email"),
		Tokens: []token{
			{
				Token:       "token1",
				Description: "description1",
				Upstream:    "upstream1",
				CreatedAt:   time.Now().Add(-time.Hour * 666).Local(),
			},
			{
				Token:       "token2",
				Description: "description2",
				Upstream:    "upstream2",
				CreatedAt:   time.Now().Add(-time.Hour * 1766).Local(),
			},
			{
				Token:       "token3",
				Description: "description3",
				Upstream:    "upstream3",
				CreatedAt:   time.Now().Add(-time.Hour * 6).Local(),
			},
		},
		Upstreams: []upstream{
			{
				Id:          "123",
				Description: "upstream1",
			},
			{
				Id:          "456",
				Description: "upstream2",
			},
		},
	}

	render := func(path string) {
		t, ok := s.templates[path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
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
