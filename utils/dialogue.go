package utils

import (
	"fmt"
	"slices"

	"github.com/eiannone/keyboard"
	"github.com/pt-main/tap/color"
)

// DialogueType represents the mode of user interaction.
type DialogueType int

const (
	// ArrowsDialogueType enables keyboard arrow navigation.
	ArrowsDialogueType DialogueType = iota
	// InputDialogueType reads a typed string from stdin.
	InputDialogueType
)

// DialogueFormat holds the prompt text and interaction mode for a dialogue.
type DialogueFormat struct {
	Text string
	Type DialogueType
}

// Run executes the dialogue and returns the selected variant.
func (df *DialogueFormat) Run(variants []string) (string, error) {
	print(color.Set(df.Text))
	if df.Type == InputDialogueType {
		return df.GetVariantByInput(variants)
	} else if df.Type == ArrowsDialogueType {
		return df.GetVariantByArrows(variants)
	} else {
		return "", fmt.Errorf("Invalid dialogue type.")
	}
}

// GetVariantByInput reads a line from stdin and validates it against allowed variants.
func (df *DialogueFormat) GetVariantByInput(variants []string) (string, error) {
	var variant string
	fmt.Scan(&variant)
	if !slices.Contains(variants, variant) {
		return "", fmt.Errorf("Invalid input: unknown variant")
	}
	return variant, nil
}

// GetVariantByArrows displays an arrow-navigable menu and returns the chosen item.
func (df *DialogueFormat) GetVariantByArrows(variants []string) (string, error) {
	if err := keyboard.Open(); err != nil {
		return "", err
	}
	defer keyboard.Close()

	var variant int
	variantsLength := len(variants)

	form := func() string {
		result := ""
		for idx, val := range variants {
			if idx == variant {
				result += "[?RD]* " + val + "[?RT]\n"
			} else {
				result += "- " + val + "\n"
			}
		}
		return color.Set(result)
	}

	rw := NewRewriter()
	rw.Write(form())

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return "", err
		}

		switch key {
		case keyboard.KeyArrowUp:
			if variant != 0 {
				variant -= 1
			}
			rw.Write(form())
		case keyboard.KeyArrowDown:
			if variant != (variantsLength - 1) {
				variant += 1
			}
			rw.Write(form())
		case keyboard.KeyEsc:
			fmt.Println("Exiting...")
			return "", err
		case keyboard.KeyEnter:
			return variants[variant], nil
		}
	}
}

// RunCheckVariant runs the dialogue and reports whether the result matches the expected variant.
func (df *DialogueFormat) RunCheckVariant(variants []string, variant string) (bool, error) {
	v, err := df.Run(variants)
	return v == variant, err
}

// NewDialogue creates a DialogueFormat with the given type and prompt text.
func NewDialogue(dtype DialogueType, form string) *DialogueFormat {
	return &DialogueFormat{
		Type: dtype,
		Text: form,
	}
}
