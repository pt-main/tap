package tap

import (
	"fmt"

	"github.com/pt-main/tap/color"
	"github.com/pt-main/tap/core"
)

func isArgsInalid(args []string, cmd command) bool {
	full_length := len(cmd.optional_args) + len(cmd.required_args)
	cond1 := (len(args) > full_length) && (!cmd.unlimited_max_args)
	cond2 := len(args) < len(cmd.required_args)
	return cond1 || cond2
}

// _call_command looks up the command by name and executes its handler.
// It validates argument count against required/optional definitions.
// Returns an error if the command is unknown or argument count is invalid.
func (p *Parser) _call_command(name string, args []string) error {
	cmd, ok := p._commands[name]
	if !ok {
		if err := p._call_subcommand(name, args); err != nil {
			return err
		}
	}
	if isArgsInalid(args, cmd) {
		return fmt.Errorf("Invalid argument length: %d.", len(args))
	}
	return cmd.handler(p, args)
}

func (p *Parser) _call_basic(args []string) error {
	cmd, ok := p._commands[DEFAULT_CMD]
	if !ok {
		return fmt.Errorf("Bad input.")
	}
	if isArgsInalid(args, cmd) {
		return fmt.Errorf("Invalid argument length: %d.", len(args))
	}
	return cmd.handler(p, args)
}

func (p *Parser) _call_subcommand(name string, args []string) error {
	cmd, ok := p._sub_commands[name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", name)
	}
	return cmd.Parse(args)
}

// _parse_args extracts flags (--flag, --key=value, --key:value) from the raw argument slice.
// Flags are stored in p.Flags (value is empty string if no value was given).
// Returns the remaining non‑flag arguments.
func (p *Parser) _parse_args(argv []string) []string {
	p.__print_verbose("Parsing args.")
	flags, res := (&core.Utils{}).ParseArgs(argv)
	p.Flags = flags
	return res
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
	println(color.Set(p._commands[DEFAULT_CMD].docstring))
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
		"Check flags by verbose and debug. \n    Flags: %v, \n    Parser flags: %v",
		p.Flags, p._parser_flags,
	)
}
