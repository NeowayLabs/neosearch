package utils

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func UniqueUint64Add(slices []uint64, i uint64) []uint64 {
	for _, e := range slices {
		if e == i {
			return slices
		}
	}

	return append(slices, i)
}

func UniqueIntAdd(slices []int, i int) []int {
	for _, e := range slices {
		if e == i {
			return slices
		}
	}

	return append(slices, i)
}
