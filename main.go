package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
	"github.com/ryodocx/private-endpoint-proxy/pkg/handlers"
)

//go:embed dist/*
var files embed.FS

func main() {
	// pprof
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		fatal(http.ListenAndServe("127.0.0.1:6060", mux))
	}()

	d, err := dao.New(nil)
	if err != nil {
		fatal(err)
	}

	// handlers
	f, err := fs.Sub(files, "dist")
	if err != nil {
		fatal(err)
	}
	h, err := handlers.New(f, d)
	if err != nil {
		fatal(err)
	}

	// TODO: graceful shutdown
	fatal(http.ListenAndServe("127.0.0.1:8080", h))
}

func fatal(v ...any) {
	fmt.Fprintln(os.Stderr, time.Now().Local())
	pc, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s:%d %s()\n", file, line, runtime.FuncForPC(pc).Name())
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
