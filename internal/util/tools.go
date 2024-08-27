package util

func Contains[T comparable](slice []T, item T) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}
