package log

import (
	"errors"
	"fmt"
)

const (
	// --- Log -----
	LogEndpointsRegistered = "%s Endpoints registered"

	// --- Error ---

	// API
	ErrBadRequest            = "Bad request - json malformed or field types incorrect"
	ErrRequiredParamNil      = "%s param cannot be empty"
	ErrInvalidAPICredentials = "Incorrect client id/api key combination"
	ErrMissingRQPrefix       = "Missing request prefix (service name and version)"
	ErrTokenExpired          = "Token invalid or expired, re-authenticate"
	ErrEndpointAuth          = "Not authorized to use endpoint"

	// Roles
	ErrRequestInvalidRole = "Invalid role name provided"
	ErrStorageNoRoles     = "No roles retrieved from storage"

	// Couchbase
	ErrCouchbaseConnection      = "Connection error"
	ErrCouchbaseAuth            = "Authentication error"
	ErrCouchbaseBucketOpen      = "Unable to open bucket"
	ErrCouchbaseSelectFailed    = "Couchbase select query failed"
	ErrCouchbaseInsertFailed    = "Couchbase insert query failed"
	ErrCouchbaseBucketRetrieval = "Bucket could not be retrieved"
	ErrCouchbaseNoResults       = "No results were returned from couchbase"
)

func GetError(str string) error {
	return errors.New(str)
}

func GetErrorf(str string, vars ...interface{}) error {
	return errors.New(fmt.Sprintf(str, vars...))
}
