package tap

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/pt-main/tap/color"
)

const Version = "0.6.0"

// Type of handler for command.
type HandlerFuncType func(*Parser, []string) error

// Inner struct for save commands to dict
type command struct {
	name               string
	handler            HandlerFuncType
	docstring          string
	required_args      []string
	optional_args      []string
	unlimited_max_args bool
}

/*
# Tap - Terminal Argument Parsing

This parser - main object in tap.

Main methods:

- AddCommand(name string, handler HandlerFuncType)

- Main() - start parser
*/
type Parser struct {
	_cli_name     string
	_about_info   string
	_parser_flags map[string]bool
	Flags         map[string]string
	_commands     map[string]command
	_config       ParserConfig
}

// Help docstring
const help_docs = "Show this message"

/*
Realization of help command.
Build autodocumentation in terminal for commans.

Docs contains name of command, args (optional and required) and description
*/
func help_cmd_handler(p *Parser, _ []string) error {
	p._print_about()
	docstrings := []string{}
	for key := range p._commands {
		el := p._commands[key]
		docs := strings.Split(el.docstring, "\n")
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
			if el.optional_args != nil || el.required_args != nil {
				color.PrintlnColored(p._config.help_args_header_block_fmt)
				args_doc := p._config.help_args_data_block_fmt
				if el.required_args != nil {
					for arg := range el.required_args {
						args_doc += "<[?RD]" + el.required_args[arg] + "[?RT]>"
						if arg != (len(el.required_args) - 1) {
							args_doc += ", "
						}
					}
				}
				if el.optional_args != nil {
					if len(args_doc) != 0 {
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
				color.PrintlnColored(args_doc)
			}
			color.PrintlnColored(p._config.help_docs_header_block_fmt)
			for line := range docs {
				color.PrintlnColored(p._config.help_docs_data_block_fmt + docs[line])
			}
			color.PrintlnColored(p._config.help_end_block_fmt)
			docstrings = append(docstrings, el.docstring)
		}
	}
	return nil
}

// Create Parser object.
func NewParser(cli_name string, about string, help_commands []string, config ParserConfig) Parser {
	p := Parser{
		_cli_name:   cli_name,
		_about_info: about,
		_parser_flags: map[string]bool{
			"debug": false, "verbose": false,
		},
		_commands: map[string]command{},
		Flags:     map[string]string{},
		_config:   config,
	}
	if help_commands == nil {
		help_commands = []string{"help", "h"}
	}
	for _, cmd := range help_commands {
		p.AddCommand(cmd, help_cmd_handler, help_docs, nil, nil, false)
	}
	return p
}

// Add command to parser.
func (p *Parser) AddCommand(
	name string,
	handler HandlerFuncType,
	docs string,
	required_args []string,
	optional_args []string,
	unlimited_max_args bool,
) {
	p.__print_verbose(
		"Adding command '%s' with %v required and %v optional args",
		name, required_args, optional_args,
	)
	p._commands[name] = command{
		name:               name,
		handler:            handler,
		docstring:          docs,
		required_args:      required_args,
		optional_args:      optional_args,
		unlimited_max_args: unlimited_max_args,
	}
}

/*
Find handler of command by name and call this.

Handler args format: [parser: Parser], [args: []string]
*/
func (p *Parser) _call_command(name string, args []string) error {
	cmd, ok := p._commands[name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", name)
	}
	full_length := len(cmd.optional_args) + len(cmd.required_args)
	cond1 := (len(args) > full_length) && (!cmd.unlimited_max_args)
	cond2 := len(args) < len(cmd.required_args)
	if cond1 || cond2 {
		return fmt.Errorf("Invalid argument length: %d.", len(args))
	}
	return cmd.handler(p, args)
}

/*
Parse args and return arguments, write flags to flags class variable.

Write "" if flag has no value (--flag=""/--flag)
*/
func (p *Parser) _parse_args(argv []string) []string {
	p.__print_verbose("Parsing args.")
	flags, res := utils{}.parse_args(argv)
	p.Flags = flags
	return res
}

// Print info if [flag] parser flag is true
func (p *Parser) Print(flag string, format string, args ...any) {
	spaces := strings.Repeat(" ", len(flag))
	format = strings.ReplaceAll(format, "\n", "\n"+spaces+" [?GN]=>[?RT] ")
	if p._parser_flags[flag] {
		color.PrintlnColored("[?RD]"+strings.ToUpper(flag)+"[?RT] [?GN]=>[?RT] "+format, args...)
	}
}

// Print info if debug parser flag is true
func (p *Parser) __print_verbose(format string, args ...any) {
	p.Print("verbose", format, args...)
}

// Print info if debug parser flag is true
func (p *Parser) __print_debug(format string, args ...any) {
	p.Print("debug", format, args...)
}

// Inner function for print about
func (p *Parser) _print_about() {
	p.__print_verbose("Print about")
	color.PrintlnColored(p._about_info)
}

// Inner function for check and enable system flags
func (p *Parser) __check_flags() {
	_, verbose_ok := p.Flags["verbose"]
	if verbose_ok {
		p._parser_flags["verbose"] = true
	}
	_, debug_ok := p.Flags["debug"]
	if debug_ok {
		p._parser_flags["debug"] = true
	}
	p.__print_verbose(
		"Check flags by verbose and debug. \nFlags: %v, Parser flags: %v",
		p.Flags, p._parser_flags,
	)
}

/*
Main function of Parser.
*/
func (p *Parser) Main() error {
	argv := p._parse_args(os.Args[1:])
	p.__check_flags()
	if len(argv) < 1 {
		p._print_about()
		var help_name string
		for _, el := range p._commands {
			if el.docstring == help_docs {
				help_name = el.name
				break
			}
		}
		color.PrintlnColored(
			"[?RD]Has no command.[?RT] Type [[?YW]%s[?RT]] for help.",
			help_name,
		)
		return errors.New("No command provided")
	}
	p.__print_verbose("Finding and calling command...")
	cmd := argv[0]
	args := argv[1:]
	p.__print_verbose(
		"Call '%s' with %v args...",
		cmd, args,
	)
	err := p._call_command(cmd, args)
	p.__print_verbose("Return after call: %v", err)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return fmt.Errorf("Command %q failed: %w", cmd, err)
	}
	return nil
}
