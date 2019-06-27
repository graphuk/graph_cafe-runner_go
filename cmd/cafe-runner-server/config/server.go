package config

type ServerConfiguration struct {
	Port     int
	Hostname string
	Cafe     struct {
		LowPort  int
		HighPort int
	}
}
