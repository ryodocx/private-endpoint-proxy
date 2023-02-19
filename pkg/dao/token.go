package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryodocx/private-endpoint-proxy/pkg/model"
)

type Token struct {
	Token       string    `db:"token"`
	UserId      string    `db:"user_id"`
	UpstreamId  string    `db:"upstream_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (t Token) Validate() (err error) {
	if t.Token == "" {
		err = errors.Join(err, fmt.Errorf("Token empty"))
	}
	if t.Description == "" {
		err = errors.Join(err, fmt.Errorf("Description empty"))
	}
	if t.UserId == "" {
		err = errors.Join(err, fmt.Errorf("UserId empty"))
	}
	if t.UpstreamId == "" {
		err = errors.Join(err, fmt.Errorf("UpstreamId empty"))
	}
	return
}

func (t Token) ToModel(upstream model.Upstream) *model.Token {
	return &model.Token{
		Token:       t.Token,
		Description: t.Description,
		Upstream:    upstream,
		CreatedAt:   t.CreatedAt,
	}
}

func NewToken(userId string, token model.Token) Token {
	return Token{
		UserId:      userId,
		Token:       uuid.NewString(),
		Description: token.Description,
		UpstreamId:  token.Upstream.Id,
		CreatedAt:   time.Now(),
	}
}
