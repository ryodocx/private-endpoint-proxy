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

	"github.com/ryodocx/private-endpoint-proxy/pkg/config"
	"github.com/ryodocx/private-endpoint-proxy/pkg/dao/dummy"
	"github.com/ryodocx/private-endpoint-proxy/pkg/handlers"
)

//go:embed dist/*
var files embed.FS

func main() {
	ifFatal := func(err error) {
		if err == nil {
			return
		}
		fmt.Fprintln(os.Stderr, time.Now().Local())
		pc, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "%s:%d %s()\n", file, line, runtime.FuncForPC(pc).Name())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// pprof
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		ifFatal(http.ListenAndServe("127.0.0.1:6060", mux))
	}()

	// config
	c, err := config.New("example/config.yaml")
	ifFatal(err)

	// dao
	d, err := dummy.New()
	// d, err := sqlite.New("tmp.db")
	ifFatal(err)

	// handlers
	f, err := fs.Sub(files, "dist")
	ifFatal(err)
	h, err := handlers.New(c, f, d)
	ifFatal(err)

	// TODO: graceful shutdown
	ifFatal(http.ListenAndServe("127.0.0.1:8080", h))
}
