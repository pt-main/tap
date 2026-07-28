package lib

import (
	"fmt"
	"strings"

	"github.com/pt-main/lc/engine"
	"github.com/pt-main/lc/engine/core"
	"github.com/pt-main/lc/parsing/stringParsing"
	"github.com/pt-main/tap/color"
)

func Color(se *engine.StringEngine, pn *stringParsing.ParsedNode) error {
	i := pn.Metadata["value"].(string)
	setter, err := core.ScopeGet[*color.Setter](se.UEP.Scope, "setter")
	if err != nil {
		return err
	}
	s := ""
	for _, col := range strings.Split(i, ";") {
		s += "[?" + col + "]"
	}
	s = setter.Set(s)
	return se.UEP.Generator.AddString(s, "main")
}

func Setup(se *engine.StringEngine, pn *stringParsing.ParsedNode) error {
	i := pn.Metadata["value"].(string)
	setter, err := core.ScopeGet[*color.Setter](se.UEP.Scope, "setter")
	if err != nil {
		return err
	}
	for _, param := range strings.Split(i, ";") {
		switch param {
		case "ColEnable", "Color+", "C+":
			*setter.ColorEnabled = true
		case "ColDisable", "Color-", "C-":
			*setter.ColorEnabled = false
		case "Col", "Color", "C":
			*setter.ColorEnabled = !(*setter.ColorEnabled)
		case "AddReset+", "AR+":
			setter.AddReset = true
		case "AddReset-", "AR-":
			setter.AddReset = false
		case "AddReset", "AR":
			setter.AddReset = !setter.AddReset
		default:
			return fmt.Errorf("Invalid param: %s", param)
		}
	}
	return nil
}

func Manual(se *engine.StringEngine, pn *stringParsing.ParsedNode) error {
	s := pn.Metadata["value"].(string)
	setter, err := core.ScopeGet[*color.Setter](se.UEP.Scope, "setter")
	if err != nil {
		return err
	}
	if !(*setter.ColorEnabled) {
		return nil
	}
	if len(s) == 0 {
		s = "0"
	}
	s = fmt.Sprintf("\033[%sm", s)
	return se.UEP.Generator.AddString(s, "main")
}

func Raw(se *engine.StringEngine, pn *stringParsing.ParsedNode) error {
	s := pn.Metadata["value"].(string)
	return se.UEP.Generator.AddString(s, "main")
}
