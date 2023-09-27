package utils

func InSlice[T comparable](slice []T, s T) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func Unique[T comparable](slice []T) []T {
	existsMap := make(map[T]bool)
	for _, v := range slice {
		existsMap[v] = true
	}
	r := make([]T, 0, len(existsMap))
	for k := range existsMap {
		r = append(r, k)
	}
	return r
}

func ReserveSlice[T any](slice []T, capacity int) []T {
	if cap(slice) < capacity {
		tmp := make([]T, len(slice), capacity)
		copy(tmp, slice)
		slice = tmp
	}
	return slice
}

func SplitSlice[T any](slice []T, n int) (ret [][]T) {
	if n <= 0 {
		panic("SplitSlice: n should be > 0")
	}
	length := len(slice)
	ret = make([][]T, 0, length/n+length%n)
	for left, right := 0, 0; left < length; left = right {
		right = left + n
		if right > length {
			right = length
		}
		ret = append(ret, slice[left:right])
	}
	return ret
}

// NOTE:non gurantee for order
func Prepend[T any](slice []T, e T) []T {
	slice = append(slice, e)
	slice[0], slice[len(slice)-1] = slice[len(slice)-1], slice[0]
	return slice
}

func StringMapEqual(lhs, rhs map[string]string) bool {
	MapEqualHelper := func(lhs, rhs map[string]string) (ret bool) {
		ret = true
		for k, v := range lhs {
			if v != rhs[k] {
				ret = false
				break
			}
		}
		return
	}
	return len(lhs) == len(rhs) && MapEqualHelper(lhs, rhs)
}
