package config

type ServerConfiguration struct {
	Port        int
	Hostname    string
	ExternalURL string
	Cafe        struct {
		LowPort  int
		HighPort int
	}
}
