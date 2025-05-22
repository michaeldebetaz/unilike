package assert

import "log"

func At[T any](slice []T, index int, msg string) T {
	if index < 0 || index >= len(slice) {
		log.Fatalf("%s: index %d out of range for slice of length %d", msg, index, len(slice))
	}

	return slice[index]
}
