package utils

import (
	"strings"

	"github.com/extemporalgenome/slug"
)

func FieldNorm(field string) string {
	fparts := strings.Split(field, ".")

	for i := range fparts {
		fparts[i] = slug.SlugAscii(fparts[i])
	}

	return strings.Join(fparts, ".")
}
