# Tap - a CLI parsing library for Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![GitHub](https://img.shields.io/badge/GitHub-repo-181717?logo=github)](https://github.com/pt-main/tap/tree/rust)
[![Docs.rs](https://docs.rs/tap-rs/badge.svg)](https://docs.rs/tap-rs)

```bash
cargo add tap-rs
# OR
cargo add --git https://github.com/pt-main/tap/rs
```

**Tap** is a lightweight library for building CLIs in Rust. Simple API for parsing arguments, commands, flags, and colored output. One dependency (`string_template`), no macros.

---

## Features

- Colored output via short codes: `[?GN]`, `[?RD]`, `[?YW]`, etc.
- Auto-generated help with grouping of commands by identical handlers.
- Subcommands (nested commands).
- Output renderers with templates and typing plus JSON serialization.

---

## Quick start

A simple CLI with a `hello` command:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    tooling::console::engine::add_help_handler,
};

fn hello_handler(_: &mut Engine, args: &[&str]) -> ErrorType {
    let name = args.first().unwrap_or(&"World");
    println!("Hello, {}!", name);
    Ok(())
}

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "MyApp v1.0".to_string());

    engine.new_command(
        "hello",
        CmdInfo::new(hello_handler, vec![], vec!["name"], false, "Say hello"),
    );

    add_help_handler(vec!["help", "-h"], &mut engine);

    if let Err(e) = engine.main() {
        eprintln!("Error: {}", e);
    }
}
```

Run:

```bash
$ cargo run hello Rust
Hello, Rust!
$ cargo run hello
Hello, World!
$ cargo run help
MyApp v1.0

╭───────  Command [help OR -h]
⎬─ Args: (none)
⎬─ Docs:
│     Generate and print help message
╰─────── 
╭───────  Command [hello]
⎬─ Args: [name]
⎬─ Docs:
│     Say hello
╰─────── 
```

The box-drawing characters in the `help` output are what `generate_help` actually prints. Not a README decoration.

---

## Commands and arguments

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType { ... }

engine.new_command(
    "copy",                               // command name
    CmdInfo::new(
        my_handler,                       // handler function
        vec!["src", "dst"],               // required arguments
        vec!["force"],                    // optional
        false,                            // unlimited: if true, accepts any number of extra arguments
        "Copy source to destination",     // description for help
    ),
);
```

The help will show:

```
╭───────  Command [copy]
⎬─ Args: <src>, <dst>, [force]
⎬─ Docs:
│     Copy source to destination
╰─────── 
```

---

## Flags

Flags are written as `--flag`, `--key=value`, `--key:value`. They are parsed automatically in `ArgsP`:

- `args_p.flags` - list of flags without values.
- `args_p.values` - dictionary key -> value.

Reading in a handler:

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType {
    if e.args_p.flags.contains(&"verbose".to_string()) {
        println!("Verbose mode enabled");
    }
    if let Some(out) = e.args_p.values.get("output") {
        println!("Output file: {}", out);
    }
    Ok(())
}
```

Important: flags are stripped from `args` before the handler is called. Access them only via `e.args_p`.

---

## Colors

Colored output via `[?COLOR]` codes:

- `[?GN]`, `[?GREEN]` - green
- `[?BGN]`, `[?BGREEN]` - bright green
- `[?BACKGREEN]`, `[?BKGN]` - green background
- `[?BACKBGREEN]`, `[?BKBGN]` - bright green background
- `[?BOLD]`, `[?BD]` - bold
- `[?DIM]`, `[?DM]` - dim
- `[?ITALIC]`, `[?IT]` - italic
- `[?UNDERLINE]`, `[?UE]` - underline
- `[?BLINK]`, `[?BL]` - blinking
- `[?REVERSE]`, `[?RV]` - reverse
- `[?STRIKETHROUGH]`, `[?ST]` - strikethrough
- `[?RT]` - reset

Format: `[?<BACKGROUND><BRIGHT><COLOR>]`.

- `<BACKGROUND>` - prefix `BK` (short) or `BACK` (full) for background.
- `<BRIGHT>` - prefix `B` for the bright variant.
- `<COLOR>` - name or two letters (first and last letter of the name). Example: `GN` = green.

### Palette

**Standard 16**: black, red, green, yellow, blue, magenta, cyan, white, and their bright variants. All are available as background.

**Extended (xterm-256)**: also available as background:

| Name      | Abbr. | Name        | Abbr. |
|-----------|-------|-------------|-------|
| orange    | `OE`  | lime        | `LE`  |
| pink      | `PK`  | teal        | `TL`  |
| purple    | `PE`  | navy        | `NY`  |
| violet    | `VT`  | gold        | `GD`  |
| brown     | `BN`  | silver      | `SR`  |
| gray      | `GY`  | maroon      | `MN`  |

Examples: `[?OE]` - orange, `[?BACKGOLD]` / `[?BKGD]` - gold background, `[?BKOE]` - orange background.

### Styles

`BOLD` (`BD`), `DIM`/`FAINT` (`DM`), `ITALIC` (`IT`), `UNDERLINE` (`UE`), `BLINK` (`BL`), `REVERSE` (`RV`), `STRIKETHROUGH` (`ST`), `RESET` (`RT`).

### Usage

```rust
use tap::tooling::colors::print::println;
use tap::tooling::colors::set::colorize;

let raw = "[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]";

println(raw);                       // prints already colored
let colored: String = colorize(raw); // or get a string
```

`println`/`print` in `colors::print` already call `colorize` internally - no need to colorize twice.

Disable colors globally:

```rust
use tap::tooling::colors::set::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

Unknown codes (`[?FOO]`) are not consumed or turned into garbage - they remain in the string as-is. This is convenient for debugging.

---

## Grouping commands in help

Commands with identical `CmdInfo` are automatically grouped into one line as aliases:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

The help will show `[help / -h]`.

`required`, `optional`, `unlimited`, and `docs` are compared. If `docs` differ, there will be no grouping.

---

## Renderers

The `tooling::renderable` module allows templating output by scope:

```rust
use tap::tooling::renderable::{
    Renderable, add_template_render, add_template_render_colored,
};

let mut r = Renderable::new();
r.scope.insert("name", "Rust");
r.scope.insert("version:int", "1");     // :int - validated as a number
r.scope.insert("stable:bool", "true");  // :bool - only "true"/"false"

add_template_render("Hello, {{name}} (v{{version}})", "plain", &mut r);
add_template_render_colored("[?GN]Hello, {{name}}[?RT]", "colored", &mut r);

println!("{}", r.render_with("plain").unwrap());
println!("{}", r.render_with("colored").unwrap());
println!("{}", r.render_json().unwrap()); // {"name":"Rust","stable":true,"version":1}
```

`render_json` returns a `Result` - for an invalid `:int`/`:bool` it returns `RenderError::InvalidValue` rather than silently inserting garbage into JSON.

---

> MIT - see the LICENSE file for details.
> 2026, By Pt.