package utils

import (
	"fmt"
	"strings"
)

func CenterText(s string, width int) string {
	runes := []rune(s)
	n := len(runes)
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

func FramedTextMap(frameColor, header, postfix string, text map[string]string) string {
	res := []string{}
	res = append(res, fmt.Sprintf("[?%v]╭─────── [?RT]%v", frameColor, header))
	for name, textSub := range text {
		if name != "" {
			res = append(res, fmt.Sprintf("[?%v]⎬─ [?RT]%v", frameColor, name))
		}
		for _, line := range strings.Split(textSub, "\n") {
			res = append(res, fmt.Sprintf("[?%v]│  [?RT]%v", frameColor, line))
		}
	}
	res = append(res, fmt.Sprintf("[?%v]╰─────── [?RT]%v", frameColor, postfix))
	return strings.Join(res, "\n")
}
