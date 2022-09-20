package utils

import "os"

func IsDebugEnabled() bool {
	debugValue := os.Getenv("DEBUG")
	if debugValue == "1" {
		return true
	}

	return false
}
