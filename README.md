# Tap - Terminal Argument Parsing

[![Go Reference](https://pkg.go.dev/badge/github.com/pt-main/tap.svg)](https://pkg.go.dev/github.com/pt-main/tap)
[![Go Report Card](https://goreportcard.com/badge/github.com/pt-main/tap)](https://goreportcard.com/report/github.com/pt-main/tap)

```bash
go get github.com/pt-main/tap
```

**Tap** is a lightweight library for building beautiful CLI applications in Go.  
It features a simple command-based API, automatic `--flag` parsing, colored output, interactive keyboard dialogs, in-place terminal rewriting, and fully customisable help messages.

## Features

- **Commands** with required / optional arguments and unlimited trailing args
- **Flag parsing** - `--flag`, `--flag=value`, `--flag:value`
- **Built‑in colour support** - shortcodes like `[?GN]`, `[?RD]`, `[?YW]` - easy and readable
- **Background colours & colour stack** - advanced ANSI control for rich layouts
- **Auto‑generated help** - groups aliases, shows arguments, respects custom format
- **Fully configurable** - change the look of `help` via `ParserConfig`, add your colors to `color.Colors`, disable colors with `color.ColorEnabled = false`
- **Hide commands** from help using `DONT_SHOW` docstring
- **Verbose / debug** flags - built‑in `--verbose` and `--debug` with conditional printing
- **Interactive dialogs** - selection menus (using input/selecting)
- **In-place terminal rewriting** - update previous output without scrolling

## Quick start

Create a simple CLI with a `hello` command:

```go
package main

import (
	"fmt"

	"github.com/pt-main/tap"
	"github.com/pt-main/tap/color"
)

func helloHandler(p *tap.Parser, args []string) error {
	color.PrintlnColored("[?GN]Hello[?RT], world! Args: %v", args)
	return nil
}

func main() {
	cfg := tap.NewParserConfig("", "", "", "", "", "") // defaults
	p := tap.NewParser("demo", "Demo CLI v1.0", nil, cfg)
	p.AddCommand("hello", helloHandler, "Prints a friendly greeting", nil, nil, true)

	if err := p.Main(); err != nil {
		fmt.Println("Error:", err)
	}
}
```

Run it:
<img alt="demo-run" src="https://github.com/user-attachments/assets/38229a2b-918e-4a43-a4f9-6e2727a2afe9" />


## Commands and arguments

Add a command using `AddCommand`:

```go
p.AddCommand(
    name               string,
    handler            HandlerFuncType, // func(*Parser, []string) error
    docstring          string,
    requiredArgs       []string,
    optionalArgs       []string,
    unlimitedMaxArgs   bool,
)
```

- **requiredArgs** - shown as `<arg>` in help. The command fails if they are missing.
- **optionalArgs** - shown as `[arg]` in help.
- **unlimitedMaxArgs** - if `true`, the command accepts any number of trailing arguments. shown as `...` in help.

You can create aliases for existing commands using `AddAlias`:

```go
p.AddCommand("help", helpHandler, tap.HELP_DOCS, nil, nil, false)
p.AddAlias("h", "help")  // now "h" works exactly like "help"
```
The alias inherits all properties (handler, arguments, docstring) from the original command.

### Example

```go
p.AddCommand("copy",
    copyHandler,
    "Copy source to destination",
    []string{"src", "dst"},  // required
    []string{"force"},       // optional
    true,
)
```

Help output would show:
```
copy <src>, <dst>, [force]...
```

## Flags

Flags are written as `--flag` or `--key=value` (also `--key:value`).  
They are parsed automatically and stored in `p.Flags` (a `map[string]string`). A flag without a value gets an empty string.

Built‑in flags:
- `--verbose` - enables verbose output (messages printed with `p.Print("verbose", ...)`)
- `--debug` - enables debug output (similar)

Your handlers can read flags directly:

```go
func myHandler(p *tap.Parser, args []string) error {
    if val, ok := p.Flags["output"]; ok {
        fmt.Println("Output file:", val)
    }
    return nil
}
```

Use `p.Print(flag, format, ...)` to output messages only when a specific flag (e.g., `--verbose` or `--debug`) is present:

```go
func myHandler(p *tap.Parser, args []string) error {
    p.Print("verbose", "Starting with args: %v", args)
    p.Print("debug", "Debug info...")
    // ...
}
```

The output is automatically prefixed with the flag name and coloured.

## Colors

<img alt="color-demo" src="https://github.com/user-attachments/assets/d5dec876-7b1d-4fcc-8ec3-9465c095a94c" />

You can just write `[?COLOR]` with uppercased color name from list to set color. Like `[?RED]` for red.

All colors: Bold, Underline, Reset, Black, Red, Green, Yellow, Blue, Magenta, Cyan.

Also you can set color using first and last letters of color. Like `[?RD]` for red.

Bright variants: `[?BRED]`, `[?BRD]`, ...

Use them with `color.PrintlnColored` or `color.PrintColored`:

```go
color.PrintlnColored("[?GN]Success[?RT] - file saved as [?YW]%s[?RT]", filename)
```

To disable colours globally:

```go
color.ColorEnabled = false
```

You can set color to string with `color.Set`:
```go
text := color.Set("[?RD]Test")
```
(Reset will be auto pasted in the end of text)

### Background colours

Background colours use the `BACK` prefix or `BK` shortcode:

```go
color.PrintlnColored("[?BKRD] ERROR: [?RT] error text...")
```

Available all colors for text except bold and underlinne.

### Colour stack

You can restore the previous colour with `[?<]` (or `[?BACK]`) and clear the stack with `[?SRT]` (or `[?SRESET]`):

```go
color.PrintlnColored("[?BE][?UE]test [?BD]bold [?RT][?<]restored")
```



## Interactive dialogs

Tap includes a small utility for interactive user prompts. It supports arrow-key navigation and validated text input.

```go
import "github.com/pt-main/tap/utils"

// Arrow-key selection
dialogue := utils.NewDialogue(utils.ArrowsDialogueType, "[?CN]Select action:[?RT]")
choice, err := dialogue.Run([]string{"install", "update", "remove"})
// Navigate with up/down, press enter to confirm, esc to cancel.
```

For simple typed input with validation against allowed variants:

```go
inputDlg := utils.NewDialogue(utils.InputDialogueType, "Enter mode (fast/slow): ")
mode, err := inputDlg.Run([]string{"fast", "slow"})
```

## Rewriting terminal output

For in-place updates (progress bars, live menus, spinners) use `Rewriter`:

```go
rw := utils.NewRewriter()
rw.Write("[?YW]Loading... 0%[?RT]")
// ... later ...
rw.Write("[?GN]Loading... 100%[?RT]")  // replaces the previous line without scrolling
```

## Customising the help output

Create a `ParserConfig` and pass it to `NewParser`.  
All fields support format strings - use `%s` for the command name or argument list.

```go
cfg := tap.NewParserConfig(
    "[?CN]>>> Command [?RT]%s[?CN] <<<[?RT]",
    "[?CN]Args:[?RT]",
    "    %s",
    "[?CN]Description:[?RT]",
    "    %s",
    "[?CN]---[?RT]",
)
p := tap.NewParser("mycli", "My tool", nil, cfg)
```

If you pass an empty string for any field, the default (coloured, nice looking) will be used.

## Grouping commands / aliases

If multiple commands share the **same docstring**, they are displayed together in help:

```go
helpDocs := "Show help"
p.AddCommand("help", helpHandler, helpDocd, nil, nil, false)
p.AddCommand("h", helpHandler, helpDocd, nil, nil, false) // Equals to p.AddAlias("h", help)
```

Help shows: `[help / h]`

## Hiding commands from help

Use `tap.DONT_SHOW` as the docstring:

```go
p.AddCommand("internal", internalHandler, tap.DONT_SHOW, nil, nil, false)
```

This command will work but will never appear in the help output.

## Full example

A minimal but complete CLI with multiple commands:

```go
package main

import (
    "fmt"
    "os"
    "github.com/pt-main/tap"
    "github.com/pt-main/tap/color"
)

func main() {
    cfg := tap.NewParserConfig("", "", "", "", "", "")
    p := tap.NewParser("myapp", "My application v0.1", nil, cfg)

    p.AddCommand("greet", func(p *tap.Parser, args []string) error {
        name := "world"
        if len(args) > 0 {
            name = args[0]
        }
        color.PrintlnColored("[?GN]Hello, %s![?RT]", name)
        return nil
    }, "Say hello", []string{}, []string{"name"}, false)

    p.AddCommand("print", func(p *tap.Parser, args []string) error {
        color.PrintlnColored("[?YW]%s[?RT]", args[0])
        return nil
    }, "Print first argument", []string{"text"}, nil, false)

    if err := p.Main(); err != nil {
        fmt.Fprintln(os.Stderr, "Fatal:", err)
        os.Exit(1)
    }
}
```

## License

MIT - see [LICENSE](LICENSE) file.  