// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/ and
// http://creativecommons.org/publicdomain/zero/1.0/legalcode

package optarg

import (
	"bytes"
	"regexp"
	"strings"
)

const (
	_ALIGN_LEFT = iota
	_ALIGN_CENTER
	_ALIGN_RIGHT
	_ALIGN_JUSTIFY
)

var reg_multilinewrap = regexp.MustCompile("[^a-zA-Z0-9,.]")

func multilineWrap(text string, linesize, leftmargin, rightmargin, alignment int) []string {
	var n int

	lines := make([]string, 0)
	pad := make([]byte, leftmargin)

	for n = 0; n < leftmargin; n++ {
		pad[n] = ' '
	}

	linesize--

	if linesize < 1 {
		linesize = 80
	}

	wordboundary := 0
	size := linesize - leftmargin - rightmargin

	if len(text) <= size {
		return []string{align(text, pad, linesize, size, alignment)}
	}

	for n = 0; n < len(text); n++ {
		if reg_multilinewrap.MatchString(text[n : n+1]) {
			wordboundary = n
		}

		if n > size {
			lines = append(lines,
				align(strings.TrimSpace(text[0:wordboundary]),
					pad, linesize, size, alignment))
			text = text[wordboundary:]
			n = 0
		}
	}

	if len(text) > 0 {
		lines = append(lines, align(strings.TrimSpace(text), pad, linesize, size, alignment))
	}

	return lines
}

func align(v string, pad []byte, linesize, size, alignment int) string {
	var buf bytes.Buffer
	switch alignment {
	case _ALIGN_LEFT:
		buf.Write(pad)
		buf.WriteString(v)

	case _ALIGN_RIGHT:
		diff := linesize - len(v) - len(pad)
		buf.Write(pad)
		for n := 0; n < diff; n++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(v)

	case _ALIGN_CENTER:
		diff := (size - len(v)) / 2
		buf.Write(pad)
		for n := 0; n < diff; n++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(v)

	case _ALIGN_JUSTIFY:
		diff := size - len(v)
		if strings.Index(v, " ") == -1 || diff == 0 {
			buf.Write(pad)
			buf.WriteString(v)
			break
		}

		var spread string
		for spread = "  "; len(v) < size; spread += " " {
			v = strings.Replace(v, spread[0:len(spread)-1], spread, -1)
		}

		for len(v) > size {
			if strings.Index(v, spread) == -1 {
				spread = spread[0 : len(spread)-1]
			}
			v = strings.Replace(v, spread, spread[0:len(spread)-1], 1)
		}

		buf.Write(pad)
		buf.WriteString(v)
	}

	return buf.String()
}
