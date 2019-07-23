package config

const (
	CurrentAPIVersion = "v1"
	ServiceName       = "gocouchbase"
	APIPrefix         = "/" + ServiceName + "/" + CurrentAPIVersion

	RequestContext = "RequestContext"

	ClientIDKey    = "clientID"
	RequestIDKey   = "requestID"
	RequestTimeKey = "timestamp"
	RemoteAddrKey  = "remoteAddr"
	MethodKey      = "method"
	RequestURIKey  = "requestURI"
)

var (
	OpenPaths = []string{"/health", "/authenticate"}
)
