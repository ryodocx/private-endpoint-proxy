package config

type config struct {
	Server  Server    `yaml:"server"`
	Auth    Auth      `yaml:"auth"`
	Backend []Backend `yaml:"backend"`
}
type Server struct {
	ListenAddr  string   `yaml:"listenAddr"`
	IPWhitelist []string `yaml:"ipWhitelist"`
}
type Headers struct {
	ID string `yaml:"id"`
}
type Proxy struct {
	Enabled bool    `yaml:"enabled"`
	Headers Headers `yaml:"headers"`
}
type Auth struct {
	Proxy Proxy `yaml:"proxy"`
}
type Backend struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
