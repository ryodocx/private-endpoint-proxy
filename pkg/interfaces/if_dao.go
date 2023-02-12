package interfaces

type Dao interface {
	GetUpstreamByToken(token string) (*Upstream, error)
	// GetTokens(userId string) ([]*Token, error)
	// CreateToken(token Token) error
	// DeleteToken(token string) error
}
