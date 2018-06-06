package fate

type NameConfig struct {
}

var defaultConfig = NewNameConfig()

func NewNameConfig() *NameConfig {
	return &NameConfig{}
}
