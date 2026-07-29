package tap

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pt-main/tap/color"
	"github.com/pt-main/tap/core"
)

// HandlerFuncType defines the signature for command handler functions.
// It receives the parser instance and the slice of command arguments.
// Returns an error if the command execution fails.
type HandlerFuncType func(*Parser, []string) error

// Command stores internal metadata for a registered command.
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

This parser is the main object in tap.

Main methods:

- AddCommand(name string, handler HandlerFuncType)

- Main() - start parser
*/
type Parser struct {
	_cli_name     string
	_about_info   string
	_parser_flags map[string]bool
	Flags         map[string]string
	Scope         core.ScopeType
	_commands     map[string]command
	_sub_commands map[string]*Parser
	_config       ParserConfig
}

// NewParser creates a new Parser instance.
// Parameters:
//   - cli_name: name of the CLI application (used in help).
//   - about: informational text printed when no command is given.
//   - help_commands: slice of command names that trigger the help handler.
//     If nil, defaults to []string{"help", "h"}.
//   - config: ParserConfig controlling help message formatting.
//
// Returns a pointer to the initialized Parser.
func NewParser(cli_name string, about string, help_commands []string, config ParserConfig) *Parser {
	p := Parser{
		_cli_name:   cli_name,
		_about_info: about,
		_parser_flags: map[string]bool{
			"debug": false, "verbose": false,
		},
		_commands:     map[string]command{},
		_sub_commands: make(map[string]*Parser),
		Flags:         map[string]string{},
		_config:       config,
	}
	if help_commands == nil {
		help_commands = []string{"help", "h"}
	}
	for _, cmd := range help_commands {
		p.AddCommand(cmd, help_cmd_handler, HELP_DOCS, nil, nil, false)
	}
	return &p
}

// AddCommand registers a new command with the parser.
// Parameters:
//   - name: command name (string used in CLI).
//   - handler: function called when the command is invoked.
//   - docs: description shown in help; use DONT_SHOW to hide the command.
//   - required_args: slice of required argument names.
//   - optional_args: slice of optional argument names.
//   - unlimited_max_args: if true, command accepts any number of trailing arguments.
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

// AddSubCommand registers a nested sub-parser for a given subcommand name.
//
// This method stores the provided Parser instance under the specified name,
// automatically integrate it into the main command dispatch flow.
//
// Parameters:
//   - name: the subcommand name.
//   - parser: the Parser instance that will handle the subcommand.
func (p *Parser) AddSubcommand(name string, parser *Parser) {
	p._sub_commands[name] = parser
}

// SubCommand retrieves a previously registered sub-parser by name.
//
// It returns the Parser pointer and a nil error if the name exists;
// otherwise, it returns an error describing the invalid name.
func (p *Parser) Subcommand(name string) (err error, pr *Parser) {
	var ok bool
	pr, ok = p._sub_commands[name]
	if !ok {
		err = fmt.Errorf("Invalid subcommand name: %v", name)
		return
	}
	return
}

// Add alias for command.
func (p *Parser) AddAlias(aliasName, cmdName string) error {
	cmdMap, ok := p._commands[cmdName]
	if !ok {
		return errors.New("[?RD]Can't add alias[?RT]: \nCommand not found")
	}
	p._commands[aliasName] = cmdMap
	return nil
}

// Print outputs a formatted message only if the given flag (e.g., "debug", "verbose") is enabled.
// The message can contain color shortcodes. Each newline is prefixed with the flag’s name for alignment.
func (p *Parser) Print(flag string, format string, args ...any) {
	spaces := strings.Repeat(" ", len(flag))
	format = strings.ReplaceAll(format, "\n", "\n"+spaces+" [?GN]=>[?RT] ")
	if p._parser_flags[flag] {
		color.PrintlnColored("[?RD]"+strings.ToUpper(flag)+"[?RT] [?GN]=>[?RT] "+format, args...)
	}
}

// Parse args.
func (p *Parser) Parse(cmdArgs []string) error {
	argv := p._parse_args(cmdArgs)
	p.__check_flags()
	if len(argv) < 1 {
		p._print_about()
		var help_name string
		for _, el := range p._commands {
			if el.docstring == HELP_DOCS {
				help_name = el.name
				break
			}
		}
		color.PrintlnColored(
			"[?RD]Has no command.[?RT] Type [[?YW]%s[?RT]] for help.",
			help_name,
		)
		return errors.New("[?RD]No command provided[?RT]")
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
		nerr := p._call_basic(argv)
		if nerr != nil {
			p.Print("debug", "'%s' cmd handler call with '%s' args end with error: %v %v", cmd, args, err, nerr)
			return fmt.Errorf("[?RD]Command [?BBK]%v[?RD] failed: \n%w", cmd, err)
		}
	}
	return nil
}

// Main is the primary entry point of the parser.
// It parses os.Args[1:], extracts flags, finds the command, and executes the corresponding handler.
// Returns an error if no command is provided, the command is unknown, or the handler fails.
func (p *Parser) Main() error {
	err := p.Parse(os.Args[1:])
	if err != nil {
		errstart := "[?RT][?YW]->[?RT]    [?BBK]|[?RT] "
		return fmt.Errorf(color.Set(
			fmt.Sprintf("[?YW]Execution stopped[?RT]: \n"+errstart+"%v",
				strings.ReplaceAll(
					err.Error(),
					"\n", "\n"+errstart,
				),
			)))
	}
	return fmt.Errorf("")
}
