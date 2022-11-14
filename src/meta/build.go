package meta

const (
	unknown string = "unknown"
)

// Application build information (injected during
// the build process)
var (
	Version    string = unknown
	GitCommit  string = unknown
	BinaryType string = unknown
	BuildTime  string = unknown
)

// GetSentByHeader returns the value of X-SentBy header
// used by default handler and Admin API
func GetSentByHeader() string {
	return "GoosyMock/" + Version
}
