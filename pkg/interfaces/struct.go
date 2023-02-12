package interfaces

import (
	"net/url"
	"time"
)

type Token struct {
	Token       string
	Description string
	UpstreamId  *Upstream
	CreatedAt   *time.Time
}

type Upstream struct {
	Id          string
	Description string
	Url         *url.URL
}
