package dao

type Dao interface {
	GetTokensByUserId(userId string) ([]*Token, error)
	GetUpstreamIdByToken(token string) (upstreamId string, ok bool, err error)
	CreateToken(token Token) error
	DeleteToken(token string) error
	MigrationUpstreamId(oldId, newId string) error
}
