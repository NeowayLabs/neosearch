package utils

func UniqueAdd(slices []uint64, i uint64) []uint64 {
	for _, e := range slices {
		if e == i {
			return slices
		}
	}

	return append(slices, i)
}
