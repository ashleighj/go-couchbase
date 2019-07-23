package couchbase

type APIClient struct {
	ClientID string `json:"client_id"`
	APIKey   string `json:"api_key"`
	Role     string `json:"role"`
}
type ServiceRole struct {
	RoleName     string   `json:"role_name"`
	AllowedPaths []string `json:"allowed_paths"`
}
