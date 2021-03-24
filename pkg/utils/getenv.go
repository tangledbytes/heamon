package utils

import "os"

// GetEnv takes name of environmental variable
// and a fallback, if the value of the environmental
// variable is "" it will return the fallback
// else it will return the value obtained
func GetEnv(env, fallback string) string {
	if res := os.Getenv(env); res != "" {
		return res
	}

	return fallback
}
