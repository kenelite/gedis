package config

type Config struct {
	ListenAddr string   // e.g., ":7379"
	Backends   []string // list of backend gedis servers
}
