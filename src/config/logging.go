package config

// LoggingConfig represents the configuration of application's logger
type LoggingConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info" json:"level"`
	Pretty bool   `env:"PRETTY_LOG" envDefault:"false" json:"pretty"`
}
