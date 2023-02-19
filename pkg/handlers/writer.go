package handlers

import (
	"log"
	"net/http"
	"strings"
)

// https://stackoverflow.com/a/43897364
type doneWriter struct {
	http.ResponseWriter
	writeDone bool
	logDone   bool
	pattern   string
	r         *http.Request
	s         *server
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
	if token := w.s.proxyPathRegex.FindStringSubmatch(w.r.URL.Path); len(token) >= 2 {
		uri = strings.Replace(uri, token[1], "<redacted>", 1)
	}
	log.Println(w.r.Method, uri, "from", w.r.RemoteAddr, w.r.Header.Get("User-Agent"), "->", w.pattern, ":", status, http.StatusText(status)) // TODO: format
}
