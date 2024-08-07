// Package utils provides utility functions for internal use.
//
// This file was generated by setup_project.sh script.
package utils

import "os"

// TODO: Implement utility functions for internal use.

const ENVKey = "RANGER_ENV"

const (
	RuntimeENVDev     = "dev"
	RuntimeENVLocal   = "local"
	RuntimeENVStaging = "staging"
	RuntimeENVProd    = "prod"
)

func Env() string {
	return os.Getenv(ENVKey)
}
