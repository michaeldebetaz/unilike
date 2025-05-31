package assert

import "log"

func At[T any](slice []T, index int) T {
	if index < 0 || index >= len(slice) {
		log.Fatalf("Index %d out of range for slice of length %d:\n%v", index, len(slice), slice)
	}

	return slice[index]
}

func NotEmpty(s string) string {
	if s == "" {
		log.Fatal("String cannot be empty")
	}

	return s
}
