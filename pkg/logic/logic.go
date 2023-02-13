package logic

import (
	"fmt"
	"net/url"

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

type Logic interface {
	GetUpstreams() ([]*model.Upstream, error)
	GetUpstreamByToken(token string) (*model.Upstream, error)
	GetTokensByUserId(userId string) ([]*model.Token, error)
}

type logic struct {
	dao dao.Dao
}

func New(dao dao.Dao) (Logic, error) {
	return &logic{
		dao: dao,
	}, nil
}

func (l *logic) GetUpstreams() ([]*model.Upstream, error) {
	u, _ := url.Parse("https://example.com")

	return []*model.Upstream{
		{
			Id:          "UpstreamID",
			Description: "Upstream Description1",
			Url:         u,
		},
		{
			Id:          "UpstreamID2",
			Description: "Upstream Description2",
			Url:         u,
		},
		{
			Id:          "UpstreamID3",
			Description: "Upstream Description3",
			Url:         u,
		},
	}, nil
}

func (l *logic) GetUpstreamByToken(token string) (*model.Upstream, error) {
	// u, _ := url.Parse("https://example.com")
	// u, _ := url.Parse("http://localhost:5555")
	u, _ := url.Parse("http://localhost:8080/")

	return &model.Upstream{
		Url: u,
	}, nil
}

func (l *logic) GetTokensByUserId(userId string) ([]*model.Token, error) {
	tokens, err := l.dao.GetTokensByUserId(userId)
	if err != nil {
		return nil, err
	}

	var modelTokens []*model.Token
	for i, t := range tokens {
		u, _ := url.Parse(fmt.Sprintf("https://example.com/%d", i))

		upstream := model.Upstream{
			Id:          fmt.Sprintf("upstream%d", i),
			Description: fmt.Sprintf("Description%d", i),
			Url:         u,
		}
		modelTokens = append(modelTokens, t.ToModel(upstream))
	}

	return modelTokens, nil
}
