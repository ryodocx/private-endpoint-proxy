package dao

import (
	"net/url"

	"github.com/ryodocx/private-endpoint-proxy/pkg/interfaces"
)

type dao struct {
	config interfaces.Config
}

func New(
	config interfaces.Config,
) (interfaces.Dao, error) {
	return &dao{
		config: config,
	}, nil
}

func (d dao) GetUpstreamByToken(token string) (*interfaces.Upstream, error) {

	u, err := url.Parse("https://example.com")
	if err != nil {
		return nil, err
	}

	return &interfaces.Upstream{
		Url: u,
	}, nil
}
