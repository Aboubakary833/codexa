package helpers

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// MakeItems is a helper function that take a slice of type T and
// a mutate function, build and return a new bubble list.Item slice.
func MakeItems[T any](s []T, mutate func(si T) list.Item) []list.Item {
	items := []list.Item{}

	for _, i := range s {
		items = append(items, mutate(i))
	}

	return items
}

// CenterHorizontally is a helper function for centering a content horizontally
func CenterHorizontally(width int, str string) string {
	return lipgloss.PlaceHorizontal(
		width, lipgloss.Center, str,
	)
}

// IsListUnfiltered check if a given is is in unfiltered state
func IsListUnfiltered(l list.Model) bool {
	return l.FilterState() == list.Unfiltered
}
