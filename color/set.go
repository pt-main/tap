package color

import (
	"strings"
)

// ColorEnabled globally enables or disables ANSI color output.
// When false, all color codes are stripped from the output.
var ColorEnabled = true

// Replace color codes
func ReplaceColors(text string) string {
	result := text
	for code, ansi := range Colors {
		if !ColorEnabled {
			ansi = ""
		}
		result = strings.ReplaceAll(result, FormColorPlaceholder(code), ansi)
	}
	return result
}

// Place color code (like 'RD', 'RED', etc.) to form for replace to color
// Like 'rd' -> '[?RD]'
func FormColorPlaceholder(code string) string {
	return "[?" + strings.ToUpper(code) + "]"
}

type Setter struct {
	ColorEnabled *bool
	AddReset     bool
}

func NewSetter(AddReset bool, ColorEnabled *bool) *Setter {
	return &Setter{
		ColorEnabled: ColorEnabled,
		AddReset:     AddReset,
	}
}

// Set replaces short color codes like [?RED] or [?GN] in the input string
// with the corresponding ANSI escape sequences. If ColorEnabled is false,
// the codes are removed entirely.
//
// Use 'BACK' (or '<') code to append last colors.
//
// Use 'SRESET' (or 'SRT') code to reset color stack.
//
// Example:
//
//	"[?BE][?UE]test [?BD]bold [?RT][?<]string" equal to "[?BE]test [?BD]bold [?RT][?BE][?UE]string"
func (s *Setter) Set(text string) string {
	result := text
	if !(*s.ColorEnabled) {
		return ReplaceColors(result)
	}
	for code, ansi := range Colors {
		result = strings.ReplaceAll(result, FormColorPlaceholder(code), ansi)
	}

	if s.AddReset {
		result += Colors["RT"]
	}
	return result
}

func Set(text string) string {
	return NewSetter(true, &ColorEnabled).Set(text)
}

// ConvertColored applies color formatting to each string in the slice
// using Set and returns a new slice with the colored strings.
func ConvertColored(text ...string) []string {
	result := []string{}
	for arg := range text {
		result = append(result, Set(text[arg]))
	}
	return result
}
