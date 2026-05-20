package tap

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/pt-main/tap/color"
)

// HandlerFuncType defines the signature for command handler functions.
// It receives the parser instance and the slice of command arguments.
// Returns an error if the command execution fails.
type HandlerFuncType func(*Parser, []string) error

// command stores internal metadata for a registered command.
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
	_commands     map[string]command
	_config       ParserConfig
}

// DONT_SHOW is a special docstring value that hides the command from the help output.
// The command remains functional but will not appear in the auto-generated help.
const DONT_SHOW = "#[DON'T SHOW]#"

// HELP_DOCS is the docstring used by the built‑in help command.
// Commands sharing this docstring will be grouped together as aliases.
const HELP_DOCS = "Generate and print help message"

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
			if el.optional_args != nil || el.required_args != nil {
				color.PrintlnColored(p._config.help_args_header_block_fmt)
				args_doc := ""
				if el.required_args != nil {
					for arg := range el.required_args {
						args_doc += "<[?RD]" + el.required_args[arg] + "[?RT]>"
						if arg != (len(el.required_args) - 1) {
							args_doc += ", "
						}
					}
				}
				if el.optional_args != nil {
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
		_commands: map[string]command{},
		Flags:     map[string]string{},
		_config:   config,
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

// _call_command looks up the command by name and executes its handler.
// It validates argument count against required/optional definitions.
// Returns an error if the command is unknown or argument count is invalid.
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

// _parse_args extracts flags (--flag, --key=value, --key:value) from the raw argument slice.
// Flags are stored in p.Flags (value is empty string if no value was given).
// Returns the remaining non‑flag arguments.
func (p *Parser) _parse_args(argv []string) []string {
	p.__print_verbose("Parsing args.")
	flags, res := utils{}.parse_args(argv)
	p.Flags = flags
	return res
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

// __print_verbose prints a formatted message when the "verbose" flag is enabled.
func (p *Parser) __print_verbose(format string, args ...any) {
	p.Print("verbose", format, args...)
}

// __print_debug prints a formatted message when the "debug" flag is enabled.
func (p *Parser) __print_debug(format string, args ...any) {
	p.Print("debug", format, args...)
}

// _print_about prints the CLI information (name/version) stored in _about_info.
func (p *Parser) _print_about() {
	p.__print_verbose("Print about")
	color.PrintlnColored(p._about_info)
}

// __check_flags enables internal verbose/debug flags based on presence in p.Flags.
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

// Main is the primary entry point of the parser.
// It parses os.Args[1:], extracts flags, finds the command, and executes the corresponding handler.
// Returns an error if no command is provided, the command is unknown, or the handler fails.
func (p *Parser) Main() error {
	argv := p._parse_args(os.Args[1:])
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
