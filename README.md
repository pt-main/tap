# Tap - Terminal Argument Parsing

[![Go Reference](https://pkg.go.dev/badge/github.com/pt-main/tap.svg)](https://pkg.go.dev/github.com/pt-main/tap)
[![Go Report Card](https://goreportcard.com/badge/github.com/pt-main/tap)](https://goreportcard.com/report/github.com/pt-main/tap)

```bash
go get github.com/pt-main/tap
```

**Tap** is a lightweight, zero-dependency library for building beautiful CLI applications in Go.  
It features a simple command-based API, automatic `--flag` parsing, colored output, and fully customisable help messages.

**Version:** `0.6.1`

## Features

- **Commands** with required / optional arguments and unlimited trailing args
- **Flag parsing** - `--flag`, `--flag=value`, `--flag:value`
- **Built‑in colour support** - shortcodes like `[?GN]`, `[?RD]`, `[?YW]` - easy and readable
- **Auto‑generated help** - groups aliases, shows arguments, respects custom format
- **Fully configurable** - change the look of `help` via `ParserConfig`
- **Hide commands** from help using `DONT_SHOW` docstring
- **Verbose / debug** flags - built‑in `--verbose` and `--debug` with conditional printing
- **Colours can be disabled** globally (`color.ColorEnabled = false`)

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

```bash
$ go build -o demo
$ ./demo hello world
Hello, world! Args: [world]
$ ./demo
Demo CLI v1.0
Has no command. Type [help] for help.
```

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
- **unlimitedMaxArgs** - if `true`, the command accepts any number of trailing arguments.

### Example

```go
p.AddCommand("copy",
    copyHandler,
    "Copy source to destination",
    []string{"src", "dst"},  // required
    []string{"force"},       // optional
    false,
)
```

Help output would show:
```
copy <src>, <dst>, [force]
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

## Colors

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
p.AddCommand("help", helpHandler, tap.HELP_DOCS, nil, nil, false)
p.AddCommand("h", helpHandler, tap.HELP_DOCS, nil, nil, false)
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
Author: Pt