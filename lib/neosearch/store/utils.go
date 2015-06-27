package store

import (
	"regexp"
	"strings"
)

func validateDatabaseName(name string) bool {
	if len(name) < 3 {
		return false
	}

	parts := strings.Split(name, ".")

	if len(parts) < 2 {
		return false
	}

	// invalid extension
	if len(parts[len(parts)-1]) < 2 {
		return false
	}

	for i := 0; i < len(parts); i++ {
		rxp := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !rxp.MatchString(parts[i]) {
			return false
		}
	}

	return true
}
