package utils

import (
	"context"
	"fmt"
	"time"
	"unicode"
)

// Function passed to [ExecWithTimeout] or [ExecWithDefaultTimeout] type
type ExecFn[T any] func(context.Context) (T, error)

// ExecWithTimeout is a utility function that pass a context with the specified duration to the callback
// function passed to it.
func ExecWithTimeout[T any](duration time.Duration, fn ExecFn[T]) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	return fn(ctx)
}

// Exec call [ExecWithTimeout] with 3 seconds as the duration
func Exec[T any](fn ExecFn[T]) (T, error) {
	return ExecWithTimeout(time.Second*3, fn)
}

// ErrorUcFirst is like [UcFirst] but for errors
func ErrorUcFirst(err error) error {
	return fmt.Errorf("%s", UcFirst(err.Error()))
}

// UcFirst transforms the first letter of a string to Uppercase
func UcFirst(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Mutate take a slice of type S and transform it into a slice of type T
func Mutate[S, T any](s []S, fn func(S) T) []T {
	var result []T
	for _, i := range s {
		result = append(result, fn(i))
	}

	return result
}
