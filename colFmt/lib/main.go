package lib

import (
	"github.com/dlclark/regexp2"
	"github.com/pt-main/lc"
	"github.com/pt-main/lc/parsing/stringParsing"
	"github.com/pt-main/lc/public"
	"github.com/pt-main/tap/color"
)

func NewLexer() *stringParsing.Lexer {
	return stringParsing.NewLexer(
		[]stringParsing.LexerRule{
			{
				Type:    "SETUP",
				Pattern: regexp2.MustCompile(`\[[sS]/(?<value>[^\]]*)\]`, 0),
			},
			{
				Type:    "COLOR",
				Pattern: regexp2.MustCompile(`\[[cC]/(?<value>[^\]]*)\]`, 0),
			},
			{
				Type:    "MANUAL",
				Pattern: regexp2.MustCompile(`\[[mM]/(?<value>[^\]]*)\]`, 0),
			},
			{
				Type:    "RAW",
				Pattern: regexp2.MustCompile(`(?<value>(?:[^[]|\[(?![sScCmM]/))+)`, 0),
			},
		},
		&stringParsing.LexerConfig{
			UseBracketBalance: false,
		},
	)
}

func NewLanguage() (*lc.EngineUniversal, error) {
	e, err := lc.NewEngineBuilder(public.StringEngineType, public.StringResType).
		WithStringParser(NewLexer()).
		WithDefaultEvents(true).
		WithPipeline([]string{"main"}).
		Build()
	if err != nil {
		return nil, err
	}
	uep, err := e.GetUEP()
	if err != nil {
		return nil, err
	}
	uep.Scope["setter"] = color.NewSetter(false, &color.ColorEnabled)
	e.NewCommandString("SETUP", Setup, "")
	e.NewCommandString("COLOR", Color, "")
	e.NewCommandString("MANUAL", Manual, "")
	e.NewCommandString("RAW", Raw, "")
	return e, nil
}
