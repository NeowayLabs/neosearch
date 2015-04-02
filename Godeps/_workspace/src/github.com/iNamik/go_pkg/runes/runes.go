package runes

// IndexRune returns the index of the first occurrence in runes of the given rune.
// It returns -1 if rune is not present in runes.
func IndexRune(runes []rune, r rune) int {
	for i, c := range runes {
		if c == r {
			return i
		}
	}
	return -1
}
