package store

import "regexp"

func validateDatabaseName(name string) bool {
	if len(name) < 3 {
		return false
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+[\.]?\.[[a-zA-Z0-9_]{2,3}$`)
	return validName.MatchString(name)
}
