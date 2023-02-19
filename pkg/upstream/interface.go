package upstream

import "github.com/ryodocx/private-endpoint-proxy/pkg/model"

type UpstreamProvider interface {
	Ping() error
	Upstreams() (map[string]*model.Upstream, error)
	Upstream(upstremId string) (upstream *model.Upstream, found bool, e error)
}
