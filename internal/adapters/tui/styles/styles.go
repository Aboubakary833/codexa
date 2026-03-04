package styles

import (
	"strings"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.AdaptiveColor{Light: "#017eb0ff", Dark: "#008cc4ff"}
	SecondaryColor = lipgloss.AdaptiveColor{Light: "#f7d88cff", Dark: "#ad7e0aff"}
	TextColor      = lipgloss.AdaptiveColor{Light: "#2e2e2e", Dark: "#dddddd"}
	DimmedColor    = lipgloss.AdaptiveColor{Light: "#555555ff", Dark: "#A49FA5"}
)

var (
	heading2Prefix = lipgloss.NewStyle().SetString(" ").
			Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true)
	heading3Prefix = lipgloss.NewStyle().SetString(" ").
			Border(lipgloss.ThickBorder(), false, false, false, true)
)

var (
	NormalTitle = lipgloss.NewStyle().
			Foreground(TextColor).
			Padding(0, 0, 0, 2).
			MaxWidth(52)

	SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).
			BorderForeground(PrimaryColor).
			Foreground(PrimaryColor).
			Padding(0, 0, 0, 1).
			MaxWidth(52)

	DimmedTitle = lipgloss.NewStyle().
			Foreground(DimmedColor).
			Padding(0, 0, 0, 2).
			MaxWidth(52)

	FilterMatch = lipgloss.NewStyle().Underline(true)

	ListTitleBarStyle = lipgloss.NewStyle().Padding(0, 0, 1, 0)

	ListFilterPromptStyle = lipgloss.NewStyle().Foreground(SecondaryColor)

	ListFilterCursorStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Width(1)

	ListTitleStyle = lipgloss.NewStyle().
			Background(PrimaryColor).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).MaxWidth(52).
			Transform(strings.ToUpper).Bold(true)
)

// GetRendererStyles return Codexa snippets markdown renderer style
func GetRendererStyles() ansi.StyleConfig {
	styleConfig := styles.PinkStyleConfig

	styleConfig.H2.StylePrimitive.Prefix = heading2Prefix.String()
	styleConfig.H3.StylePrimitive.Prefix = heading3Prefix.String()

	styleConfig.HorizontalRule.Color = strPtr("#dddddd")
	styleConfig.Heading.Color = strPtr("#028abfff")
	styleConfig.Link.Color = strPtr("#f3ba2aff")

	styleConfig.Code = styles.TokyoNightStyleConfig.Code
	styleConfig.CodeBlock = styles.TokyoNightStyleConfig.CodeBlock

	return styleConfig
}

// strPtr take a string str and return its pointer
func strPtr(str string) *string {
	return &str
}
