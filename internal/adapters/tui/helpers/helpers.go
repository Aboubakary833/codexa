package helpers

import (
	"github.com/aboubakary833/codexa/utils"
	"github.com/charmbracelet/bubbles/list"
)

// MakeItems is a helper function that take a slice of type T and
// a mutate function, build and return a new bubble list.Item slice.
func MakeItems[T any](s []T, mutate func(si T) list.Item) []list.Item {
	return utils.Mutate(s, mutate)
}

// IsListUnfiltered check if a given is is in unfiltered state
func IsListUnfiltered(l list.Model) bool {
	return l.FilterState() == list.Unfiltered
}
