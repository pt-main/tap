package tap

import (
	"slices"
	"sort"
	"strings"

	"github.com/pt-main/tap/go/color"
)

func helpFormArgs(el command) string {
	parts := make([]string, 0, len(el.required_args)+len(el.optional_args)+1)
	for _, arg := range el.required_args {
		parts = append(parts, "<[?RD]"+arg+"[?RT]>")
	}
	for _, arg := range el.optional_args {
		parts = append(parts, "[[?BE]"+arg+"[?RT]]")
	}
	if el.unlimited_max_args {
		parts = append(parts, "...")
	}
	return strings.Join(parts, ", ")
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printCommandBlock(p *Parser, el command, gatherAliases bool) {
	var header string
	if gatherAliases {
		aliases := make([]string, 0, 1)
		for _, n := range sortedKeys(p._commands) {
			c := p._commands[n]
			if c.docstring == el.docstring && c.name != DEFAULT_CMD {
				aliases = append(aliases, c.name)
			}
		}
		header = "[?YW]" + strings.Join(aliases, " [?RT]/[?YW] ") + "[?RT]"
	} else {
		header = "[?YW]" + el.name + "[?RT]"
	}
	color.PrintlnColored(p._config.help_command_block_fmt, header)

	if argsDoc := helpFormArgs(el); argsDoc != "" {
		color.PrintlnColored(p._config.help_args_header_block_fmt)
		color.PrintlnColored(p._config.help_args_data_block_fmt, argsDoc)
	}

	if el.docstring != "" && el.docstring != DONT_SHOW {
		color.PrintlnColored(p._config.help_docs_header_block_fmt)
		for _, line := range strings.Split(el.docstring, "\n") {
			color.PrintlnColored(p._config.help_docs_data_block_fmt, line)
		}
	}
	color.PrintlnColored(p._config.help_end_block_fmt)
}

func printSubcommandBlock(p *Parser, name string, el *Parser) {
	color.PrintlnColored(p._config.help_subcommand_block_fmt, "[?YW]"+name+"[?RT]")
	if el._about_info != "" {
		color.PrintlnColored(p._config.help_docs_header_block_fmt)
		for _, line := range strings.Split(el._about_info, "\n") {
			color.PrintlnColored(p._config.help_docs_data_block_fmt, line)
		}
	}
	color.PrintlnColored(p._config.help_end_block_fmt)
}

func help_cmd_handler(p *Parser, args []string) error {
	if len(args) > 0 {
		return helpForTargets(p, args)
	}
	p._print_about()
	helpPrintAll(p)
	return nil
}

func helpPrintAll(p *Parser) {
	shownDocstrings := make([]string, 0, len(p._commands))

	for _, name := range sortedKeys(p._commands) {
		el := p._commands[name]

		if el.docstring == DONT_SHOW || el.name == DEFAULT_CMD {
			shownDocstrings = append(shownDocstrings, el.docstring)
			continue
		}
		if slices.Contains(shownDocstrings, el.docstring) {
			continue
		}
		printCommandBlock(p, el, true)
		shownDocstrings = append(shownDocstrings, el.docstring)
	}

	for _, name := range sortedKeys(p._sub_commands) {
		printSubcommandBlock(p, name, p._sub_commands[name])
	}
}

func helpForTargets(p *Parser, targets []string) error {
	for _, query := range targets {
		if el, ok := p._commands[query]; ok {
			printCommandBlock(p, el, true)
			continue
		}
		if el, ok := p._sub_commands[query]; ok {
			printSubcommandBlock(p, query, el)
			continue
		}
		color.PrintlnColored("[?RD]Unknown command:[?RT] [?YW]%s[?RT]", query)
	}
	return nil
}
