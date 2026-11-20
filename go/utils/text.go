package utils

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pt-main/tap/go/color"
	"golang.org/x/term"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func visibleLen(s string) int {
	cleaned := ansiRegex.ReplaceAllString(s, "")
	return utf8.RuneCountInString(cleaned)
}

func skipANSI(runes []rune, idx int) (int, bool) {
	if runes[idx] == '\x1b' && idx+1 < len(runes) && runes[idx+1] == '[' {
		for j := idx + 2; j < len(runes); j++ {
			if (runes[j] >= 'a' && runes[j] <= 'z') || (runes[j] >= 'A' && runes[j] <= 'Z') {
				return j, true
			}
		}
	}
	return idx, false
}

func CenterText(s string, width int) string {
	n := visibleLen(s)
	if n >= width {
		return s
	}
	left := (width - n) / 2
	right := width - n - left
	return fmt.Sprintf("%*s%s%*s", left, "", s, right, "")
}

func FramedTextList(frameColor, header, postfix string, text []string) string {
	textMap := map[string]string{}
	for _, key := range text {
		textMap[key] = ""
	}
	return FramedTextMap(frameColor, header, postfix, textMap)
}

const FramedTextDefaultTermWidthPadding = 10

func FramedTextMap(frameColor, header, postfix string, text map[string]string) string {
	return FramedTextMapWithPadding(frameColor, header, postfix, FramedTextDefaultTermWidthPadding, text)
}

func FramedTextMapWithPadding(frameColor, header, postfix string, padding int, text map[string]string) string {
	res := []string{}
	res = append(res, fmt.Sprintf("[?%v]╭─────── [?RT]%v", frameColor, header))
	for _, name := range slices.Sorted(maps.Keys(text)) {
		textSub := text[name]
		if name != "" {
			res = append(res, fmt.Sprintf("[?%v]⎬─ [?RT]%v", frameColor, name))
		}
		textSub = CutTextForWidth(max(GetTermWidth()-8-padding, 10), color.Set(textSub))
		for _, line := range strings.Split(textSub, "\n") {
			res = append(res, fmt.Sprintf("[?%v]│  [?RT]%v", frameColor, line))
		}
	}
	res = append(res, fmt.Sprintf("[?%v]╰─────── [?RT]%v", frameColor, postfix))
	return strings.Join(res, "\n")
}

// returns -1 if has error
func GetTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = -1
	}
	return width
}

func isVowel(r rune) bool {
	const vowels = "aeiouyаеёиоуыэюя"
	return strings.ContainsRune(vowels, unicode.ToLower(r))
}

func isConsonant(r rune) bool {
	r = unicode.ToLower(r)
	return unicode.IsLetter(r) && !isVowel(r)
}

func CutTextForWidth(width int, text string) string {
	if width <= 0 {
		return text
	}

	var sb strings.Builder
	sb.Grow(len(text) + len(text)/width)

	runes := []rune(text)
	lineLen := 0

	for i := 0; i < len(runes); i++ {
		if endIdx, ok := skipANSI(runes, i); ok {
			for k := i; k <= endIdx; k++ {
				sb.WriteRune(runes[k])
			}
			i = endIdx
			continue
		}

		char := runes[i]

		if char == '\n' {
			sb.WriteRune(char)
			lineLen = 0
			continue
		}

		if lineLen >= width {
			var nextChar rune
			nextIdx := i + 1
			for nextIdx < len(runes) {
				if endIdx, ok := skipANSI(runes, nextIdx); ok {
					nextIdx = endIdx + 1
					continue
				}
				nextChar = runes[nextIdx]
				break
			}

			if char == ' ' {
				sb.WriteRune('\n')
				lineLen = 0
				continue
			}

			sb.WriteRune(char)

			if !unicode.IsLetter(char) {
				if nextChar == ' ' {
					for k := i + 1; k < nextIdx; k++ {
						sb.WriteRune(runes[k])
					}
					i = nextIdx
				}
				sb.WriteRune('\n')
				lineLen = 0
			} else if nextChar == ' ' || nextChar == 0 || !unicode.IsLetter(nextChar) {
				if nextChar == ' ' {
					for k := i + 1; k < nextIdx; k++ {
						sb.WriteRune(runes[k])
					}
					i = nextIdx
				}
				sb.WriteRune('\n')
				lineLen = 0
			} else if isVowel(char) || isConsonant(char) {
				sb.WriteString("-\n")
				lineLen = 0
			}
			continue
		}

		sb.WriteRune(char)
		lineLen++
	}

	return sb.String()
}
