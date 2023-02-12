package config

import (
	"io/ioutil"
	"log"

	"github.com/goccy/go-yaml"
	"github.com/ryodocx/private-endpoint-proxy/pkg/util"
)

func init() {
	b, err := ioutil.ReadFile("example/config.yaml")
	if err != nil {
		util.Fatal(err)
	}
	log.Println("\n" + string(b))

	v := config{}
	if err := yaml.Unmarshal(b, &v); err != nil {
		util.Fatal(err)

	}
	log.Println(v)
}
