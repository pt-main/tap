package tap

import (
	"slices"
	"strings"

	"github.com/pt-main/tap/color"
)

func help_form_args(el command) string {
	if el.optional_args != nil || el.required_args != nil || el.unlimited_max_args {
		args_doc := ""
		if el.required_args != nil && len(el.required_args) != 0 {
			for arg := range el.required_args {
				args_doc += "<[?RD]" + el.required_args[arg] + "[?RT]>"
				if arg != (len(el.required_args) - 1) {
					args_doc += ", "
				}
			}
		}
		if el.optional_args != nil && len(el.optional_args) != 0 {
			if len(args_doc) > 2 {
				args_doc += ", "
			}
			for arg := range el.optional_args {
				args_doc += "[[?BE]" + el.optional_args[arg] + "[?RT]]"
				if arg != (len(el.optional_args) - 1) {
					args_doc += ", "
				}
			}
		}
		if el.unlimited_max_args {
			args_doc += "..."
		}
		return args_doc
	}
	return ""
}

// help_cmd_handler implements the built‑in help command.
// It prints a formatted help message listing all visible commands,
// their arguments (required/optional), and descriptions.
func help_cmd_handler(p *Parser, _ []string) error {
	p._print_about()
	docstrings := []string{}
	for key := range p._commands {
		el := p._commands[key]
		docs := strings.Split(el.docstring, "\n")
		if el.docstring == DONT_SHOW {
			docstrings = append(docstrings, el.docstring)
		}
		if slices.Index(docstrings, el.docstring) == -1 {
			cmds := []string{}
			for key := range p._commands {
				if p._commands[key].docstring == el.docstring {
					cmds = append(cmds, p._commands[key].name)
				}
			}
			commands := "[?YW]"
			for idx, cmd := range cmds {
				commands += cmd
				if idx != (len(cmds) - 1) {
					commands += " [?RT]/[?YW] "
				}
			}
			commands += "[?RT]"
			color.PrintlnColored(p._config.help_command_block_fmt, commands)
			args_doc := help_form_args(el)
			if args_doc != "" {
				color.PrintlnColored(p._config.help_args_header_block_fmt)
				color.PrintlnColored(p._config.help_args_data_block_fmt, args_doc)
			}
			color.PrintlnColored(p._config.help_docs_header_block_fmt)
			for line := range docs {
				color.PrintlnColored(p._config.help_docs_data_block_fmt, docs[line])
			}
			color.PrintlnColored(p._config.help_end_block_fmt)
			docstrings = append(docstrings, el.docstring)
		}
	}
	return nil
}
