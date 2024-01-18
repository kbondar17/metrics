package app

type ServerConfig struct {
	Address *string
}

type AppConfig struct {
	Server ServerConfig
}

func NewAppConfig(host *string) *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Address: host,
		},
	}
}
