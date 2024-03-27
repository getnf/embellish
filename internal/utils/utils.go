package utils

func Filter[T any](items []T, condition func(T) bool) (results []T) {
	for _, item := range items {
		if condition(item) {
			results = append(results, item)
		}
	}
	return
}
