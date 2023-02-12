package dao

import (
	"net/url"

	"github.com/ryodocx/private-endpoint-proxy/pkg/interfaces"
)

type dao struct {
	upstreams []*interfaces.Upstream
}

func New() (interfaces.Dao, error) {
	return &dao{}, nil
}

func (d dao) GetUpstreamByToken(token string) (*interfaces.Upstream, error) {

	u, err := url.Parse("http://localhost:8080/livez")
	if err != nil {
		return nil, err
	}

	return &interfaces.Upstream{
		Url: u,
	}, nil
}
