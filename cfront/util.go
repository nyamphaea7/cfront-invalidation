package cfront

func contain[T comparable](arr []T, v T) bool {
	for _, vInArr := range arr {
		if vInArr == v {
			return true
		}
	}
	return false
}
