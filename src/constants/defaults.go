package constants

import "net/http"

// Configuration defaults
const (
	DefaultCfgAdminAPIAddr        string = ":8081"
	DefaultCfgAdminAPISecuredAddr string = ":8444"
	DefaultCfgContentAddr         string = ":8080"
	DefaultCfgContentSecuredAddr  string = ":8443"
)

// Response defaults
const (
	DefaultResponseStatusCode  int    = http.StatusOK
	DefaultResponseContentType string = "text/plain"
)
