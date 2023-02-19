package static

import (
	"net/url"

	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
	"github.com/ryodocx/private-endpoint-proxy/pkg/upstream"
)

type static struct {
	upstream map[string]*model.Upstream
}

func (s static) Ping() error {
	return nil
}

func (s static) Upstreams() (map[string]*model.Upstream, error) {
	return s.upstream, nil
}

func (s static) Upstream(upstremId string) (upstream *model.Upstream, found bool, e error) {
	upstreams, err := s.Upstreams()
	if err != nil {
		return nil, false, err
	}
	u, ok := upstreams[upstremId]
	return u, ok, nil
}

func New(id, description, upstreamUrl []string) (upstream.UpstreamProvider, error) {
	upstream := map[string]*model.Upstream{}
	for i, id := range id {
		u, err := url.Parse(upstreamUrl[i])
		if err != nil {
			return nil, err
		}
		upstream[id] = &model.Upstream{
			Id:          id,
			Description: description[i],
			Url:         u,
		}
	}
	return &static{
		upstream: upstream,
	}, nil
}
