package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
)

type sqlite struct {
	db *sql.DB
}

func New(filepath string) (dao.Dao, error) {

	db, err := sqlx.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}
	fmt.Println(db.Exec(`
  CREATE TABLE IF NOT EXISTS tokens(
    token       INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT
  );
  INSERT INTO test (name) VALUES('aaa');
  INSERT INTO test (name) VALUES('aaa');
  INSERT INTO test (name) VALUES('bbb');
  INSERT INTO test (name) VALUES('aaa');
  `))

	// db.Get()

	return &sqlite{
		db: db.DB,
	}, nil
}

func (v *sqlite) GetTokensByUserId(userId string) ([]*dao.Token, error) {
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

func (v *sqlite) GetUpstreamIdByToken(token string) (upstreamId string, ok bool, err error) {
	return "upstream1", true, nil
}

func (v *sqlite) CreateToken(token dao.Token) error {
	return nil
}

func (v *sqlite) DeleteToken(token string) error {
	return nil
}

func (v *sqlite) MigrationUpstreamId(oldId, newId string) error {
	return nil
}
