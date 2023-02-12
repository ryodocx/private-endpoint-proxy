package main

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/pprof"

	"github.com/ryodocx/private-endpoint-proxy/pkg/server"
	"github.com/ryodocx/private-endpoint-proxy/pkg/util"
)

//go:embed dist/*
var files embed.FS

func main() {
	// pprof
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		util.Fatal(http.ListenAndServe("127.0.0.1:6060", mux))
	}()

	f, err := fs.Sub(files, "dist")
	if err != nil {
		util.Fatal(err)
	}

	mux, err := server.New(f)
	if err != nil {
		util.Fatal(err)
	}

	util.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
