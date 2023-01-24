package utils

func Bound[T any](list []T, limit uint) []T {
	if uint(len(list)) > limit {
		return list[:limit]
	}
	return list
}
