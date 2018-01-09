package util

func Sub(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)

	for _, val := range slice2 {
		m[val] = true
	}

	r := make([]int64, 0)
	for _, val := range slice1 {
		if !m[val] {
			r = append(r, val)
		}
	}
	return r
}
