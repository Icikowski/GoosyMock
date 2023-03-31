package config

// Config represents the application's configuration
type Config struct {
	Logging            LoggingConfig `json:"logging"`
	AdminAPIService    ServiceConfig `envPrefix:"ADMIN_API_" json:"adminApi"`
	ContentService     ServiceConfig `envPrefix:"CONTENT_" json:"content"`
	HealthProbesAddr   string        `env:"HEALTH_PROBES_ADDR" envDefault:":8888" json:"healthProbesAddr"`
	MaximumPayloadSize int64         `env:"MAX_PAYLOAD_SIZE" envDefault:"64" json:"maxPayloadSize"`
}
