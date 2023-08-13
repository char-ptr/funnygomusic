package utils

func GetIndex[T any](slice []T, idx int) T {
	for k, v := range slice {
		if k == idx {
			return v
		}
	}
	var res T
	return res
}
