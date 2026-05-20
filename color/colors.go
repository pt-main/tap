// Package color provides ANSI color formatting for terminal output.
// It supports short color codes like [?GN] (green), [?RD] (red), [?BOLD], etc.
// Coloring can be globally disabled by setting color.ColorEnabled = false.
package color

// Color codes for ANSI escape sequences.
const (
	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	brightBlack   = "\033[90m"
	brightRed     = "\033[91m"
	brightGreen   = "\033[92m"
	brightYellow  = "\033[93m"
	brightBlue    = "\033[94m"
	brightMagenta = "\033[95m"
	brightCyan    = "\033[96m"
	brightWhite   = "\033[97m"

	reset     = "\033[0m"
	bold      = "\033[1m"
	underline = "\033[4m"
)

// Colors maps short color codes (uppercase) to ANSI escape sequences.
// Supported codes: BOLD/BD, UNDERLINE/UE, RESET/RT,
// basic colors (BLACK/BK, RED/RD, GREEN/GN, YELLOW/YW, BLUE/BE, MAGENTA/MA, CYAN/CN),
// and their bright variants (BBLACK/BBK, BRED/BRD, BGREEN/BGN, BYELLOW/BYW, BBLUE/BBE, BMAGENTA/BMA, BCYAN/BCN).
var Colors = map[string]string{
	"BOLD": bold, "BD": bold,
	"UNDERLINE": underline, "UE": underline,
	"RESET": reset, "RT": reset,

	"BLACK": black, "BK": black,
	"RED": red, "RD": red,
	"GREEN": green, "GN": green,
	"YELLOW": yellow, "YW": yellow,
	"BLUE": blue, "BE": blue,
	"MAGENTA": magenta, "MA": magenta,
	"CYAN": cyan, "CN": cyan,

	"BBLACK": brightBlack, "BBK": brightBlack,
	"BRED": brightRed, "BRD": brightRed,
	"BGREEN": brightGreen, "BGN": brightGreen,
	"BYELLOW": brightYellow, "BYW": brightYellow,
	"BBLUE": brightBlue, "BBE": brightBlue,
	"BMAGENTA": brightMagenta, "BMA": brightMagenta,
	"BCYAN": brightCyan, "BCN": brightCyan,
}
