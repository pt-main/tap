package color

import (
	"fmt"
	"strings"
)

// ColorEnabled globally enables or disables ANSI color output.
// When false, all color codes are stripped from the output.
var ColorEnabled = true

// Set replaces short color codes like [?RED] or [?GN] in the input string
// with the corresponding ANSI escape sequences. If ColorEnabled is false,
// the codes are removed entirely. The ANSI reset code is automatically
// appended at the end when colors are enabled.
func Set(text string) string {
	result := text
	for code, ansi := range Colors {
		if !ColorEnabled {
			ansi = ""
		}
		result = strings.ReplaceAll(result, "[?"+code+"]", ansi)
	}
	if !ColorEnabled {
		return result
	}
	return result + Colors["RESET"]
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

// PrintColored formats the string using fmt.Sprintf, replaces color codes,
// and writes the result to stdout without a trailing newline.
func PrintColored(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Print(Set(formatted))
}

// PrintlnColored formats the string using fmt.Sprintf, replaces color codes,
// and writes the result to stdout followed by a newline.
func PrintlnColored(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Println(Set(formatted))
}
