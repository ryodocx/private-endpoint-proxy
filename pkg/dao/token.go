package dao

import (
	"time"

	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

type Token struct {
	Token       string
	Description string
	UserId      string
	UpstreamId  string
	CreatedAt   time.Time
}

func (t Token) ToModel(upstream model.Upstream) *model.Token {
	return &model.Token{
		Token:       t.Token,
		Description: t.Description,
		Upstream:    upstream,
		CreatedAt:   t.CreatedAt,
	}
}
