package config

import (
	"io/ioutil"
	"net/url"

	"github.com/goccy/go-yaml"
	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

type Config interface {
	Upstreams() (map[string]*model.Upstream, error)
	Upstream(upstremId string) (upstream *model.Upstream, found bool, e error)
}

type config struct {
	Server  Server    `yaml:"server"`
	Auth    Auth      `yaml:"auth"`
	Backend []Backend `yaml:"backend"`
}
type Server struct {
	ListenAddr  string   `yaml:"listenAddr"`
	IPWhitelist []string `yaml:"ipWhitelist"`
}
type Headers struct {
	ID string `yaml:"id"`
}
type Proxy struct {
	Enabled bool    `yaml:"enabled"`
	Headers Headers `yaml:"headers"`
}
type Auth struct {
	Proxy Proxy `yaml:"proxy"`
}
type Backend struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func (c config) Upstreams() (map[string]*model.Upstream, error) {
	// TODO
	u, _ := url.Parse("http://localhost:8080/livez")
	return map[string]*model.Upstream{
		"Id": {
			Id:          "Id",
			Description: "Description",
			Url:         u,
		},
	}, nil
}

func (c config) Upstream(upstremId string) (upstream *model.Upstream, found bool, e error) {
	upstreams, err := c.Upstreams()
	if err != nil {
		return nil, false, err
	}
	u, ok := upstreams[upstremId]
	return u, ok, nil
}

func New(filePath string) (Config, error) {

	// TODO
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// log.Println("\n" + string(b))
	v := config{}
	if err := yaml.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	// log.Println(v)

	return &config{}, nil
}
