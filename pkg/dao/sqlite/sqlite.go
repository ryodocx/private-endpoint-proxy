package sqlite

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/ryodocx/private-endpoint-proxy/pkg/dao"
)

type sqlite struct {
	db                       *sqlx.DB
	stmtCreate               *sqlx.NamedStmt
	stmtGetTokensByUserId    *sqlx.Stmt
	stmtGetUpstreamIdByToken *sqlx.Stmt
	stmtDeleteToken          *sqlx.Stmt
}

func New(filepath string) (dao.Dao, error) {
	db, err := sqlx.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS tokens(
        token       TEXT PRIMARY KEY,
        user_id     TEXT NOT NULL,
        upstream_id TEXT NOT NULL,
        description TEXT NOT NULL,
        created_at  DATETIME NOT NULL
    );
    CREATE INDEX IF NOT EXISTS tokens_user_id ON tokens(user_id);
    CREATE INDEX IF NOT EXISTS tokens_upstream_id ON tokens(upstream_id);
    `); err != nil {
		return nil, err
	}

	stmtCreate, err := db.PrepareNamed(`
    INSERT INTO tokens (
      token,
      user_id,
      upstream_id,
      description,
      created_at
    )
    VALUES (
      :token,
      :user_id,
      :upstream_id,
      :description,
      :created_at
    );
    `)
	if err != nil {
		return nil, err
	}

	stmtGetTokensByUserId, err := db.Preparex(`
    SELECT
      token,
      user_id,
      upstream_id,
      description,
      created_at
    FROM tokens
    WHERE user_id = $1;
    `)
	if err != nil {
		return nil, err
	}

	stmtGetUpstreamIdByToken, err := db.Preparex(`
    SELECT upstream_id
    FROM tokens
    WHERE token = $1;
    `)
	if err != nil {
		return nil, err
	}

	stmtDeleteToken, err := db.Preparex(`
    DELETE
    FROM tokens
    WHERE token = $1;
    `)
	if err != nil {
		return nil, err
	}

	return &sqlite{
		db:                       db,
		stmtCreate:               stmtCreate,
		stmtGetTokensByUserId:    stmtGetTokensByUserId,
		stmtGetUpstreamIdByToken: stmtGetUpstreamIdByToken,
		stmtDeleteToken:          stmtDeleteToken,
	}, nil
}

func (v *sqlite) Ping() error {
	return v.db.Ping()
}

func (v *sqlite) GetTokensByUserId(userId string) (tokens []*dao.Token, err error) {
	err = v.stmtGetTokensByUserId.Select(&tokens, userId)
	return
}

func (v *sqlite) GetUpstreamIdByToken(token string) (upstreamId string, ok bool, err error) {
	rows, err := v.stmtGetUpstreamIdByToken.Queryx(token)
	if err != nil {
		return "", false, err
	}
	defer rows.Close()
	for {
		if !rows.Next() {
			break
		}
		rows.Scan(&upstreamId)
	}
	if upstreamId != "" {
		ok = true
	}
	return
}

func (v *sqlite) CreateToken(token dao.Token) error {
	_, err := v.stmtCreate.Exec(token)
	return err
}

func (v *sqlite) DeleteToken(token string) error {
	_, err := v.stmtDeleteToken.Exec(token)
	log.Println(err)
	return err
}
