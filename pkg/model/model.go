package model

import (
	"net/url"
	"time"
)

type Context struct {
	User      string
	Tokens    []*Token
	Upstreams []*Upstream
}

type Token struct {
	Token       string
	Description string
	Upstream    Upstream
	CreatedAt   time.Time
}

type Upstream struct {
	Id          string
	Description string
	Url         *url.URL
}
