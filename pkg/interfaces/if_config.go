package interfaces

type Config interface {
	GetUpstream(id string) (*Upstream, error)
}
