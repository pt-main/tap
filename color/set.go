package color

import (
	"slices"
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

// Set replaces short color codes like [?RED] or [?GN] in the input string
// with the corresponding ANSI escape sequences. If ColorEnabled is false,
// the codes are removed entirely. The ANSI reset code is automatically
// appended at the end when colors are enabled.
//
// Use 'BACK' (or '<') code to append last colors.
//
// Use 'SRESET' (or 'SRT') code to reset color stack.
//
// Example:
//
//	"[?BE][?UE]test [?BD]bold [?RT][?<]string" equal to "[?BE]test [?BD]bold [?RT][?BE][?UE]string"
func Set(text string) string {
	BackVariants := []string{"BACK", "<"}
	SrtVariants := []string{"SRESET", "SRT"}
	result := text
	if !ColorEnabled {
		return ReplaceColors(result)
	}
	colorStack := []string{}
	for code, ansi := range Colors {
		if slices.Contains(BackVariants, code) {
			replace := ""
			for _, code := range colorStack {
				replace += Colors[code]
			}
			result = strings.Replace(result, code, replace, 1)
		} else if slices.Contains(SrtVariants, code) {
			result = strings.Replace(result, code, "", 1)
			colorStack = []string{}
		} else {
			result = strings.ReplaceAll(result, FormColorPlaceholder(code), ansi)
			if ansi == Colors["RT"] {
				if len(colorStack)-1 >= 0 {
					colorStack = colorStack[:len(colorStack)-1]
				}
			} else {
				colorStack = append(colorStack, code)
			}
		}

	}
	return result + Colors["RT"]
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
