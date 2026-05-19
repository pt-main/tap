package tap

// Configuration of Parser
type ParserConfig struct {
	help_command_block_fmt     string
	help_args_header_block_fmt string
	help_args_data_block_fmt   string
	help_docs_header_block_fmt string
	help_docs_data_block_fmt   string
	help_end_block_fmt         string
}

// Create config for parser
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
		help_args_data_block_fmt = "[?GN]│[?RT]    "
	}
	if help_docs_header_block_fmt == "" {
		help_docs_header_block_fmt = "[?GN]⎬─ Desc:[?RT]"
	}
	if help_docs_data_block_fmt == "" {
		help_docs_data_block_fmt = "[?GN]│[?RT]     "
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
