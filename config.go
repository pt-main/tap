// Package tap provides a lightweight command-line argument parser with support
// for commands, flags, colored output, and customizable help messages.
package tap

// DONT_SHOW is a special docstring value that hides the command from the help output.
// The command remains functional but will not appear in the auto-generated help.
const DONT_SHOW = "#[DON'T SHOW]#"

// DEFAULT_CMD is a reserved command name.
// It is used when no command is specified as the first argument,
// allowing your CLI to be invoked in two ways: cli args... (implicit default)
// or cli cmd args... (explicit command).
// Note: documentation of default command is displayed after about message.
const DEFAULT_CMD = "#[DEFAULT CMD]#"

// HELP_DOCS is the docstring used by the built‑in help command.
// Commands sharing this docstring will be grouped together as aliases.
const HELP_DOCS = "Generate and print help message"

// ParserConfig defines the formatting templates for the auto-generated help message.
// Each field is a format string that may contain "%s" placeholders for dynamic content.
// If an empty string is passed to NewParserConfig, the corresponding default format will be used.
type ParserConfig struct {
	help_command_block_fmt     string
	help_args_header_block_fmt string
	help_args_data_block_fmt   string
	help_docs_header_block_fmt string
	help_docs_data_block_fmt   string
	help_end_block_fmt         string
}

// NewParserConfig creates a new ParserConfig with the given format strings.
// Any empty string parameter will be replaced with a sensible default.
// Parameters:
//   - help_command_block_fmt: format for the command name block (e.g., "╭─────── Command [%s]").
//   - help_args_header_block_fmt: format for the arguments section header.
//   - help_args_data_block_fmt: format for each argument line.
//   - help_docs_header_block_fmt: format for the description section header.
//   - help_docs_data_block_fmt: format for each line of the description.
//   - help_end_block_fmt: format for the closing block.
//
// Returns a populated ParserConfig.
func NewParserConfig(
	help_command_block_fmt string,
	help_args_header_block_fmt string,
	help_args_data_block_fmt string,
	help_docs_header_block_fmt string,
	help_docs_data_block_fmt string,
	help_end_block_fmt string,
) ParserConfig {
	if help_command_block_fmt == "" {
		help_command_block_fmt = "[?GN]╭─────── Command[?RT] [%s]"
	}
	if help_args_header_block_fmt == "" {
		help_args_header_block_fmt = "[?GN]⎬─ Args:[?RT]"
	}
	if help_args_data_block_fmt == "" {
		help_args_data_block_fmt = "[?GN]│[?RT]     %s"
	}
	if help_docs_header_block_fmt == "" {
		help_docs_header_block_fmt = "[?GN]⎬─ Desc:[?RT]"
	}
	if help_docs_data_block_fmt == "" {
		help_docs_data_block_fmt = "[?GN]│[?RT]     %s"
	}
	if help_end_block_fmt == "" {
		help_end_block_fmt = "[?GN]╰───────[?RT]"
	}

	return ParserConfig{
		help_command_block_fmt:     help_command_block_fmt,
		help_args_header_block_fmt: help_args_header_block_fmt,
		help_args_data_block_fmt:   help_args_data_block_fmt,
		help_docs_header_block_fmt: help_docs_header_block_fmt,
		help_docs_data_block_fmt:   help_docs_data_block_fmt,
		help_end_block_fmt:         help_end_block_fmt,
	}
}
