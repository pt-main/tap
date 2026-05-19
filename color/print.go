package color

import (
	"fmt"
	"strings"
)

/*
Set applies colors to text using shortcodes like [?RED], [?GREEN], [?BOLD]

Also you can set color using shortcuts - first and last letters of color
to set color like [?RD], [?GN], [?BD].
*/
func Set(text string) string {
	result := text
	for code, ansi := range Colors {
		result = strings.ReplaceAll(result, "[?"+code+"]", ansi)
	}
	return result + Colors["RESET"]
}

// Convert list of strings to colored strings using Set
func ConvertColored(text ...string) []string {
	result := []string{}
	for arg := range text {
		result = append(result, Set(text[arg]))
	}
	return result
}

// PrintColored prints formatted text with color codes replaced
func PrintColored(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Print(Set(formatted))
}

// PrintlnColored prints text with color codes replaced and adds newline
func PrintlnColored(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Println(Set(formatted))
}
