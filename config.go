package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/goccy/go-yaml"
)

type config struct {
	Server    Server     `yaml:"server"`
	Database  Database   `yaml:"database"`
	Upstreams []Upstream `yaml:"upstreams"`
}

type Server struct {
	ListenAddr string `yaml:"listenAddr"`
}

type Database struct {
	SQLite SQLite `yaml:"sqlite"`
}

type SQLite struct {
	Filepath string `yaml:"filepath"`
}

type Upstream struct {
	Id          string `yaml:"id"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

func loadConfig(filePath string) (*config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	c := &config{}
	if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(b))), c); err != nil {
		return nil, err
	}

	b2, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b2))

	return c, nil
}
