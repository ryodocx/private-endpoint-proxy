package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao/sqlite"
	"github.com/ryodocx/private-endpoint-proxy/pkg/handlers"
	"github.com/ryodocx/private-endpoint-proxy/pkg/upstream/static"
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

	// flag
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	// config
	c, err := loadConfig(configPath)
	ifFatal(err)

	// pprof
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		ifFatal(http.ListenAndServe("127.0.0.1:6060", mux))
	}()

	// upstream provider
	u, err := static.New(
		func() (id, description, upstreamUrl []string) {
			for _, u := range c.Upstreams {
				id = append(id, u.Id)
				description = append(description, u.Description)
				upstreamUrl = append(upstreamUrl, u.URL)
			}
			return id, description, upstreamUrl
		}())
	ifFatal(err)

	// dao
	d, err := sqlite.New(c.Database.SQLite.Filepath)
	ifFatal(err)

	// handlers
	f, err := fs.Sub(files, "dist")
	ifFatal(err)
	h, err := handlers.New(f, u, d)
	ifFatal(err)

	// TODO: graceful shutdown
	ifFatal(http.ListenAndServe(c.Server.ListenAddr, h))
}
