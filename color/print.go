package color

import "fmt"

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
