package main

import (
	"fmt"

	"github.com/pt-main/tap/utils"
)

func main() {
	d := utils.NewDialogue(utils.ArrowsDialogueType, "")
	fmt.Println(d.GetVariantByArrows([]string{"1", "2"}))
}
