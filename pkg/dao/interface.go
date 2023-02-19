package dao

type Dao interface {
	Ping() error
	GetTokensByUserId(userId string) ([]*Token, error)
	GetUpstreamIdByToken(token string) (upstreamId string, found bool, err error)
	CreateToken(token Token) error
	DeleteToken(token string) error
	// MigrateUpstreamId(oldId, newId string) error
}
