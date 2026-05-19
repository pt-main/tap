package tap

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pt-main/tap/color"
)

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
	flags         map[string]string
	_commands     map[string]command
}

/*
Realization of help command.
Build autodocumentation in terminal for commans.

Docs contains name of command, args (optional and required) and description
*/
func help_cmd_handler(p *Parser, _ []string) error {
	p._print_about()
	for key := range p._commands {
		el := p._commands[key]
		docs := strings.Split(el.docstring, "\n")
		color.PrintlnColored("[?GN]╭─────── Command[?RT] [[?YW]%s[?RT]]", el.name)
		if el.optional_args != nil || el.required_args != nil {
			color.PrintlnColored("[?GN]⎬─ Args:[?RT]")
			args_doc := "[?GN]│[?RT]    "
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
		color.PrintlnColored("[?GN]⎬─ Desc:[?RT]")
		for line := range docs {
			color.PrintlnColored("[?GN]│[?RT]    %s", docs[line])
		}
		color.PrintlnColored("[?GN]╰───────[?RT]")
	}
	return nil
}

// Create Parser object.
func NewParser(cli_name string, about string) Parser {
	p := Parser{
		_cli_name:   cli_name,
		_about_info: about,
		_parser_flags: map[string]bool{
			"debug": false, "verbose": false,
		},
		_commands: map[string]command{},
		flags:     map[string]string{},
	}
	p.AddCommand("help", help_cmd_handler, "Show this message", nil, nil, false)
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
		return fmt.Errorf("Key not in dict: %s", name)
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
	p.flags = flags
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
	_, verbose_ok := p.flags["verbose"]
	if verbose_ok {
		p._parser_flags["verbose"] = true
	}
	_, debug_ok := p.flags["debug"]
	if debug_ok {
		p._parser_flags["debug"] = true
	}
	p.__print_verbose(
		"Check flags by verbose and debug. \nFlags: %v, Parser flags: %v",
		p.flags, p._parser_flags,
	)
}

/*
Main function of Parser.
*/
func (p *Parser) Main() error {
	argv := p._parse_args(os.Args[1:])
	p.__check_flags()
	if len(argv) < 1 {
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
