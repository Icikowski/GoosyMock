package config

// Config represents the application's configuration
type Config struct {
	Logging            LoggingConfig `json:"logging"`
	AdminAPIService    ServiceConfig `envPrefix:"ADMIN_API_" json:"adminApi"`
	ContentService     ServiceConfig `envPrefix:"CONTENT_" json:"content"`
	HealthProbesPort   int           `env:"HEALTH_PROBES_PORT" envDefault:"8888" json:"healthProbesPort"`
	MaximumPayloadSize int64         `env:"MAX_PAYLOAD_SIZE" envDefault:"64" json:"maxPayloadSize"`
}
