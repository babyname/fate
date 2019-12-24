package fate

type DBConfig struct {
}

type Config struct {
	Database DBConfig
}

func DefaultConfig() *Config {
	return &Config{}
}
