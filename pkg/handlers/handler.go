package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/ryodocx/private-endpoint-proxy/pkg/logic"
)

type server struct {
	mux           *http.ServeMux
	files         http.FileSystem
	templates     map[string]*template.Template
	consolePrefix string
	logic         logic.Logic
	regex         *regexp.Regexp
}

func New(files fs.FS, logic logic.Logic) (http.Handler, error) {
	s := &server{
		mux:           http.NewServeMux(),
		files:         http.FS(files),
		templates:     map[string]*template.Template{},
		consolePrefix: "/console/",
		logic:         logic,
		regex:         regexp.MustCompile(`^/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})(/|$|\?)`),
	}

	// load templates
	var dirwalk func(f fs.FS, dir string) (paths []string, err error)
	dirwalk = func(f fs.FS, dir string) (paths []string, err error) {
		entries, err := fs.ReadDir(f, dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				nestedPaths, err := dirwalk(f, path.Join(dir, entry.Name()))
				if err != nil {
					return nil, err
				}
				paths = append(paths, nestedPaths...)
				continue
			}
			paths = append(paths, path.Join(dir, entry.Name()))
		}
		return paths, nil
	}
	paths, err := dirwalk(files, ".")
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed at call dirwalk()"))
	}
	for _, p := range paths {
		if !strings.HasSuffix(p, ".html") {
			continue
		}
		t, err := template.ParseFS(files, p)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed at call dirwalk()"))
		}
		s.templates[p] = t
	}

	// setup handlers
	handleFunc := func(pattern string, handers ...http.HandlerFunc) {
		// TODO: support sentry.io
		s.mux.HandleFunc(
			pattern,
			func(handlers ...http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					dw := &doneWriter{
						s:              s,
						pattern:        pattern,
						r:              r,
						ResponseWriter: w,
					}
					defer func() {
						recover()
						if !dw.writeDone {
							dw.WriteHeader(http.StatusInternalServerError)
						}
					}()
					for _, h := range handlers {
						h.ServeHTTP(dw, r)
						if dw.writeDone {
							return
						}
					}
				}
			}(handers...),
		)
	}

	handleFunc("/livez",
		method(http.MethodGet, http.MethodHead),
		s.live,
	)
	handleFunc("/readyz",
		method(http.MethodGet, http.MethodHead),
		s.ready,
	)
	handleFunc(s.consolePrefix,
		method(http.MethodGet),
		s.console,
	)
	handleFunc("/api/token/add",
		method(http.MethodPost),
		antiCSRF(),
		s.addToken,
	)
	handleFunc("/api/token/delete",
		method(http.MethodPost),
		antiCSRF(),
		s.deleteToken,
	)
	handleFunc("/",
		s.proxy,
	)

	return s, nil
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// https://stackoverflow.com/a/43897364
type doneWriter struct {
	writeDone bool
	http.ResponseWriter
	logDone bool
	pattern string
	r       *http.Request
	s       *server
}

func (w *doneWriter) WriteHeader(status int) {
	w.writeDone = true
	w.ResponseWriter.WriteHeader(status)
	w.logging(status)
}

func (w *doneWriter) Write(b []byte) (int, error) {
	w.writeDone = true
	i, err := w.ResponseWriter.Write(b)
	if err != nil {
		w.logging(-1)
	} else {
		w.logging(http.StatusOK)
	}
	return i, err
}

func (w *doneWriter) logging(status int) {
	if w.logDone {
		return
	}
	w.logDone = true
	uri := w.r.RequestURI
	if token := w.s.regex.FindStringSubmatch(w.r.URL.Path); len(token) >= 2 {
		uri = strings.Replace(uri, token[1], "<redacted>", 1)
	}
	log.Println(w.r.Method, uri, "from", w.r.RemoteAddr, w.r.Header.Get("User-Agent"), "->", w.pattern, ":", status, http.StatusText(status)) // TODO: format
}
