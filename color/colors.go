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

	// Background colours (standard)
	blackBg   = "\033[40m"
	redBg     = "\033[41m"
	greenBg   = "\033[42m"
	yellowBg  = "\033[43m"
	blueBg    = "\033[44m"
	magentaBg = "\033[45m"
	cyanBg    = "\033[46m"
	whiteBg   = "\033[47m"

	// Background colours (bright)
	brightBlackBg   = "\033[100m"
	brightRedBg     = "\033[101m"
	brightGreenBg   = "\033[102m"
	brightYellowBg  = "\033[103m"
	brightBlueBg    = "\033[104m"
	brightMagentaBg = "\033[105m"
	brightCyanBg    = "\033[106m"
	brightWhiteBg   = "\033[107m"

	reset     = "\033[0m"
	bold      = "\033[1m"
	underline = "\033[4m"
)

// Colors maps short color codes (uppercase) to ANSI escape sequences.
// Supported codes:
//   - Text foreground: BOLD/BD, UNDERLINE/UE, RESET/RT,
//     basic colours (BLACK/BK, RED/RD, GREEN/GN, YELLOW/YW, BLUE/BE, MAGENTA/MA, CYAN/CN),
//     bright variants (BBLACK/BBK, BRED/BRD, BGREEN/BGN, BYELLOW/BYW, BBLUE/BBE, BMAGENTA/BMA, BCYAN/BCN)
//   - Background colours (full name: BACK<COLOUR>, short: BK<SHORT>):
//     BACKBLACK/BKBK, BACKRED/BKRD, BACKGREEN/BKGN, BACKYELLOW/BKYW, BACKBLUE/BKBE, BACKMAGENTA/BKMA, BACKCYAN/BKCN,
//     and bright variants: BACKBBLACK/BKBBK, BACKBRED/BKBRD, BACKBGREEN/BKBGN, BACKBYELLOW/BKBYW, BACKBBLUE/BKBBE, BACKBMAGENTA/BKBMA, BACKBCYAN/BKBCN.
var Colors = map[string]string{
	"BOLD": bold, "BD": bold,
	"UNDERLINE": underline, "UE": underline,
	"RESET": reset, "RT": reset,

	// Text foreground
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

	// Background colours
	"BACKBLACK": blackBg, "BKBK": blackBg,
	"BACKRED": redBg, "BKRD": redBg,
	"BACKGREEN": greenBg, "BKGN": greenBg,
	"BACKYELLOW": yellowBg, "BKYW": yellowBg,
	"BACKBLUE": blueBg, "BKBE": blueBg,
	"BACKMAGENTA": magentaBg, "BKMA": magentaBg,
	"BACKCYAN": cyanBg, "BKCN": cyanBg,

	"BACKBBLACK": brightBlackBg, "BKBBK": brightBlackBg,
	"BACKBRED": brightRedBg, "BKBRD": brightRedBg,
	"BACKBGREEN": brightGreenBg, "BKBGN": brightGreenBg,
	"BACKBYELLOW": brightYellowBg, "BKBYW": brightYellowBg,
	"BACKBBLUE": brightBlueBg, "BKBBE": brightBlueBg,
	"BACKBMAGENTA": brightMagentaBg, "BKBMA": brightMagentaBg,
	"BACKBCYAN": brightCyanBg, "BKBCN": brightCyanBg,

	// Other
	"BACK": "", "<": "",
	"SRESET": "", "SRT": "",
}
