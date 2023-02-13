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

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao/dummy"
	"github.com/ryodocx/private-endpoint-proxy/pkg/handlers"
	"github.com/ryodocx/private-endpoint-proxy/pkg/logic"
)

//go:embed dist/*
var files embed.FS

func main() {
	// pprof
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		ifFatal(http.ListenAndServe("127.0.0.1:6060", mux))
	}()

	// dao
	d, err := dummy.New()
	ifFatal(err)

	// logic
	l, err := logic.New(d)
	ifFatal(err)

	// handlers
	f, err := fs.Sub(files, "dist")
	ifFatal(err)
	h, err := handlers.New(f, l)
	ifFatal(err)

	// TODO: graceful shutdown
	ifFatal(http.ListenAndServe("127.0.0.1:8080", h))
}

func ifFatal(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, time.Now().Local())
	pc, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s:%d %s()\n", file, line, runtime.FuncForPC(pc).Name())
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
