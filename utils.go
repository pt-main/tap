package tap

import "strings"

// Class with utils for Parser
type utils struct{}

/*
Parse args for flags and args.

Read arguments like 'arg', '--flag' and '--flag=key' (or '--flag:key')
*/
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
