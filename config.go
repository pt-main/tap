// Package tap provides a lightweight command-line argument parser with support
// for commands, flags, colored output, and customizable help messages.
package tap

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
