package utils

import (
	"os"
	"strconv"
)

// GetEnvParam : return string environmental param if exists, otherwise return default
func GetEnvParam(param string, dflt string) string {
	if v, exists := os.LookupEnv(param); exists {
		return v
	}
	return dflt
}

// GetEnvBoolParam : return bool environmental param if exists, otherwise return default
func GetEnvBoolParam(param string, dflt bool) bool {
	if v, exists := os.LookupEnv(param); exists {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return dflt
		}
		return b
	}
	return dflt
}
