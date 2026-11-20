package main

import (
	"fmt"

	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/utils"
)

func main() {
	p := tap.NewParser("test", utils.FramedTextMap("RD", "", "", map[string]string{"": `[?RD]TEST TAP PARSER[?RT]
very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long description`}), []string{"help", "-h"}, tap.DefaultParserConfig())

	p.AddCommand(tap.DEFAULT_CMD, func(p *tap.Parser, s []string) error {
		fmt.Println(p.Flags, s)
		fmt.Println(p.Super.RawArgs, p.RawArgs)
		return nil
	}, "---", nil, nil, true)

	p.AddCommand("1", func(p *tap.Parser, s []string) error {
		fmt.Println("1")
		return nil
	}, "print 1", nil, nil, false)
	p.AddAlias("p1", "1")

	p.AddCommand("2", func(p *tap.Parser, s []string) error {
		fmt.Println("2")
		return nil
	}, "[?UE]very very[?RT] very very very very very very very very very very very very very very very very very very very very very very very very very very [?BD]very very[?RT] long docs",
		nil, nil, false)
	p.AddAlias("p2", "2")
	p.AddSubcommand("self", p)

	if err := p.Main(); err != nil {
		fmt.Println(err)
	}
}
