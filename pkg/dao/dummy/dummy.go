package dummy

import (
	"time"

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
)

type dummy struct{}

func New() (dao.Dao, error) {
	return &dummy{}, nil
}

func (d *dummy) Ping() error {
	return nil
}

func (d *dummy) GetTokensByUserId(userId string) ([]*dao.Token, error) {
	return []*dao.Token{
		{
			Token:       "Token1",
			Description: "Description1",
			UpstreamId:  "Upstream1",
			CreatedAt:   time.Now().Add(-time.Hour * 400),
		},
		{
			Token:       "Token2",
			Description: "Description2",
			UpstreamId:  "Upstream2",
			CreatedAt:   time.Now().Add(-time.Hour * 400),
		},
		{
			Token:       "Token3",
			Description: "Description3",
			UpstreamId:  "Upstream3",
			CreatedAt:   time.Now().Add(-time.Hour * 400),
		},
	}, nil
}

func (d *dummy) GetUpstreamIdByToken(token string) (upstreamId string, ok bool, err error) {
	return "upstream1", true, nil
}

func (d *dummy) CreateToken(token dao.Token) error {
	return nil
}

func (d *dummy) DeleteToken(token string) error {
	return nil
}

func (d *dummy) MigrationUpstreamId(oldId, newId string) error {
	return nil
}
