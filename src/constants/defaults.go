package constants

import "net/http"

// Configuration defaults
const (
	DefaultCfgAdminAPIPort        int = 8081
	DefaultCfgAdminAPISecuredPort int = 8444
	DefaultCfgContentPort         int = 8080
	DefaultCfgContentSecuredPort  int = 8443
)

// Response defaults
const (
	DefaultResponseStatusCode  int    = http.StatusOK
	DefaultResponseContentType string = "text/plain"
)
