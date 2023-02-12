package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/ryodocx/private-endpoint-proxy/pkg/util"
)

type server struct {
	mux           *http.ServeMux
	files         http.FileSystem
	templates     map[string]*template.Template
	consolePrefix string
}

func New(files fs.FS) (http.Handler, error) {
	s := &server{
		mux:           http.NewServeMux(),
		files:         http.FS(files),
		templates:     map[string]*template.Template{},
		consolePrefix: "/console/",
	}

	// load templates
	paths, err := util.Dirwalk(files, ".")
	if err != nil {
		util.Fatal(err)
	}
	for _, p := range paths {
		if !strings.HasSuffix(p, ".html") {
			continue
		}
		t, err := template.ParseFS(files, p)
		if err != nil {
			util.Fatal(err)
		}
		s.templates[p] = t
	}

	// setup handlers
	addHandleFunc := func(pattern string, handers ...http.HandlerFunc) {
		s.mux.HandleFunc(
			pattern,
			func(handlers ...http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					dw := &doneWriter{
						pattern:        pattern,
						r:              r,
						ResponseWriter: w,
					}
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

	addHandleFunc("/livez",
		method(http.MethodGet, http.MethodHead),
		s.live,
	)
	addHandleFunc("/readyz",
		method(http.MethodGet, http.MethodHead),
		s.ready,
	)
	addHandleFunc(s.consolePrefix,
		method(http.MethodGet),
		s.console,
	)
	addHandleFunc("/api/token/add",
		method(http.MethodPost),
		antiCSRF(),
		s.addToken,
	)
	addHandleFunc("/api/token/delete",
		method(http.MethodPost),
		antiCSRF(),
		s.deleteToken,
	)
	addHandleFunc("/",
		method(http.MethodGet),
		s.proxy,
	)

	return s, nil
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// https://stackoverflow.com/a/43897364
type doneWriter struct {
	http.ResponseWriter
	writeDone bool
	pattern   string
	r         *http.Request
	logDone   bool
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
	log.Println(w.r.Method, w.r.RequestURI, "->", w.pattern, ":", status, http.StatusText(status), "from", w.r.RemoteAddr, w.r.Header.Get("User-Agent")) // TODO
	w.logDone = true
}
