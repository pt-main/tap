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
	lines := strings.Split(text, "\n")

	for lineIdx, line := range lines {
		if lineIdx > 0 {
			sb.WriteByte('\n')
		}

		words := strings.Split(line, " ")
		curLen := 0

		for _, word := range words {
			wordLen := visibleLen(word)
			if wordLen > width {
				if curLen > 0 {
					sb.WriteByte('\n')
					curLen = 0
				}
				runes := []rune(word)
				chunkLen := 0
				for i := 0; i < len(runes); i++ {
					if endIdx, ok := skipANSI(runes, i); ok {
						for k := i; k <= endIdx; k++ {
							sb.WriteRune(runes[k])
						}
						i = endIdx
						continue
					}
					if chunkLen == width-1 && i < len(runes)-1 {
						sb.WriteString("-\n")
						chunkLen = 0
					}
					sb.WriteRune(runes[i])
					chunkLen++
				}
				curLen = chunkLen
				continue
			}

			needSpace := 0
			if curLen > 0 {
				needSpace = 1
			}

			if curLen+needSpace+wordLen > width {
				sb.WriteByte('\n')
				curLen = 0
				needSpace = 0
			}

			if needSpace > 0 {
				sb.WriteByte(' ')
				curLen++
			}
			sb.WriteString(word)
			curLen += wordLen
		}
	}

	return sb.String()
}
