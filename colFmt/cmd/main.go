package main

import (
	"fmt"
	"os"

	"github.com/pt-main/colFmt/lib"
	"github.com/pt-main/lc/engine/core"
)

func main() {
	if a := os.Args; len(a) > 1 {
		printErr := func(err error) {
			fmt.Printf("{%s}{Err=%s}", a[1], err.Error())
		}
		l, err := lib.NewLanguage()
		if err != nil {
			printErr(err)
			return
		}
		err = l.ProcessString(a[1])
		if err != nil {
			printErr(err)
			return
		}
		uep, _ := l.GetUEP()
		res, err := core.GetStringRes(uep.Generator, "")
		if err != nil {
			printErr(err)
			return
		}
		fmt.Print(res)
	}
}
