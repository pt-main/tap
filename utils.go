package tap

import "strings"

// utils is an internal helper type providing argument parsing functionality.
type utils struct{}

// parse_args processes a slice of command-line arguments.
// It detects flags prefixed with "--" (e.g., "--flag", "--key=value", "--key:value").
// Flags without a value are stored with an empty string.
// Returns a map of flag names to their values, and a slice of non-flag arguments.
func (u utils) parse_args(argv []string) (map[string]string, []string) {
	flags := make(map[string]string)
	var result []string

	for _, el := range argv {
		if strings.HasPrefix(el, "--") {
			el = el[2:]

			var key, value string
			var hasValue bool

			if strings.Contains(el, "=") {
				parts := strings.SplitN(el, "=", 2)
				key, value = parts[0], parts[1]
				hasValue = true
			} else if strings.Contains(el, ":") {
				parts := strings.SplitN(el, ":", 2)
				key, value = parts[0], parts[1]
				hasValue = true
			} else {
				key = el
				hasValue = false
			}

			if hasValue {
				flags[key] = value
			} else {
				flags[key] = ""
			}
		} else {
			result = append(result, el)
		}
	}

	return flags, result
}
