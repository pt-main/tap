# Tap - a CLI parsing library for Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![GitHub](https://img.shields.io/badge/GitHub-repo-181717?logo=github)](https://github.com/pt-main/tap/tree/rust)
[![Docs.rs](https://docs.rs/tap-rs/badge.svg)](https://docs.rs/tap-rs)

```bash
cargo add tap-rs
```

**Tap** is a small library for building CLIs in Rust. It parses arguments, commands, and flags, prints colored text, and generates help. One dependency: `string_template`; macros are optional.

## Features

- Colored output via short codes `[?GN]`, `[?RD]`, `[?YW]`.
- Automatic help generation with grouping of commands that share the same handler.
- Subcommands.
- Typed parsing of positional arguments via `ArgsExt`.
- Flags `--flag`, `--key=value`, `--key:value` in `ArgsP`.
- Macros `commands!`, `subcommands!`, `args!`, `bail!`, `ensure!`, `cprintln!`, and others.
- Renderers with templates, typing, and JSON serialization.

## Quick start

A simple CLI with a `hello` command:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    tooling::add_help_handler,
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

The frame characters in the `help` output are what `generate_help` actually prints. They are not README decoration.

## Quick start with macros

The same thing, but shorter:

```rust
use tap::{
    argsparsing::ArgsP, commands,
    engine::structs::{Engine, ErrorType},
    tooling::add_help_handler,
};

fn hello_handler(_: &mut Engine, args: &[&str]) -> ErrorType {
    let name = args.first().copied().unwrap_or("World");
    println!("Hello, {}!", name);
    Ok(())
}

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "MyApp v1.0".to_string());

    commands! {
        engine;
        "hello" => hello_handler { optional: ["name"], docs: "Say hello" }
    }

    add_help_handler(vec!["help", "-h"], &mut engine);

    if let Err(e) = engine.main() {
        eprintln!("Error: {}", e);
    }
}
```

The `commands!` macro creates `CmdInfo` and registers the command itself. See the Macros section below.

## Commands and arguments

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType { ... }

engine.new_command(
    "copy",
    CmdInfo::new(
        my_handler,                    // handler
        vec!["src", "dst"],            // required
        vec!["force"],                 // optional
        false,                         // unlimited: true - accepts any number of arguments
        "Copy source to destination",  // description for help
    ),
);
```

Help will show:

```
╭───────  Command [copy]
⎬─ Args: <src>, <dst>, [force]
⎬─ Docs:
│     Copy source to destination
╰─────── 
```

### Subcommands

Subcommands are commands that receive only arguments, without `&mut Engine`. Their handler has the type `fn(&[&str]) -> ErrorType`.

```rust
fn echo_handler(raw: &[&str]) -> ErrorType {
    println!("echo: {}", raw.join(" "));
    Ok(())
}

engine.new_subcommand(
    "echo",
    CmdInfo::new(echo_handler, vec![], vec![], true, "Echo back args"),
);
```

Register via `new_subcommand` / `new_subcommands` or the `subcommands!` macro.

### Parsing arguments via `ArgsExt`

The `ArgsExt` trait provides typed access to a slice of arguments:

```rust
use tap::argsparsing::ArgsExt;

fn add_handler(_: &mut Engine, raw: &[&str]) -> ErrorType {
    let a: f64 = raw.arg(0)?;
    let b: f64 = raw.arg(1)?;
    let label = raw.arg_or(2, "sum");   // if the argument is absent, returns "sum"
    println!("{}: {} + {} = {}", label, a, b, a + b);
    Ok(())
}
```

Errors contain the argument index, not a bare parser error.

## Flags

Flags are written as `--flag`, `--key=value`, or `--key:value`. They are parsed automatically in `ArgsP`:

- `args_p.flags` - list of flags without values.
- `args_p.values` - dictionary key -> value.

Flags are removed from `args` before the handler is called. Inspect them only via `e.args_p` or the `Engine` helper methods.

```rust
fn my_handler(e: &mut Engine, _: &[&str]) -> ErrorType {
    if e.flag("verbose") {
        println!("verbose mode");
    }
    if let Some(out) = e.value("output") {
        println!("output file: {}", out);
    }
    let workers: u16 = e.value_t("workers")?.unwrap_or(4);
    println!("workers: {}", workers);
    Ok(())
}
```

- `Engine::flag(name)` - `bool`.
- `Engine::value(name)` - `Option<&str>`.
- `Engine::value_t::<T>(name)` - `Result<Option<T>, String>`; a parse error comes with the flag name.

### `--` as a separator

`--` disables flag parsing. Everything after it goes into `args` as is:

```
myapp run -- --x --y    # args = ["--x", "--y"], flags = []
```

## Macros

All macros are exported from the crate root, i.e. they are called as `tap::commands!`, `tap::args!`, etc.

### `commands!` and `subcommands!`

```rust
tap::commands! {
    engine;
    "greet" => greet_handler {
        optional: ["name"],
        docs: "Greet someone",
    }
    "add" | "sum" => add_handler {
        required: ["a", "b"],
        docs: "Sum two numbers",
    }
    "echo" => echo_handler {
        unlimited: true,
        docs: "Echo args",
    }
}
```

The `required`, `optional`, and `unlimited` fields are optional; `docs` is mandatory. `"name" | "alias"` registers multiple names with one handler.

`subcommands!` works the same way, but the handler must be `fn(&[&str]) -> ErrorType`.

### `args!`

Binds positional arguments to typed variables:

```rust
fn add_handler(_: &mut Engine, raw: &[&str]) -> ErrorType {
    tap::args! { raw;
        a: f64,
        b: f64,
        label: String = "sum".to_string(),
    }
    println!("{}: {} + {} = {}", label, a, b, a + b);
    Ok(())
}
```

The form `name: Type` is a required argument; `name: Type = default` is optional. Each argument consumes one position; order matters. An error comes with the argument name and index.

### `bail!` and `ensure!`

```rust
fn h(_: &mut Engine, raw: &[&str]) -> ErrorType {
    tap::ensure!(!raw.is_empty(), "at least one argument required");
    let n: u32 = raw[0].parse().map_err(|_| format!("bad number: {}", raw[0]))?;
    tap::ensure!(n > 0, "must be positive, got {}", n);
    Ok(())
}
```

`bail!` returns `Err` early from the handler. `ensure!` does the same, but conditionally.

### `cprint!`, `cprintln!`, `ceprintln!`

Analogues of `print!` / `println!` / `eprintln!`, but with automatic coloring of `[?COLOR]` codes:

```rust
tap::cprintln!("[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]");
```

### `scope!`, `template!`, `template_colored!`

Shorthands for working with `Renderable`:

```rust
use tap::tooling::renderable::Renderable;

let mut r = Renderable::new();
tap::scope!(r;
    "name"        => "Rust",
    "version:int" => "1",
    "stable:bool" => "true",
);
tap::template!(r, "plain", "Hello, {{name}} (v{{version}})");
tap::template_colored!(r, "colored", "[?GN]Hello, {{name}}[?RT]");
```

## Colors

Colored output via `[?COLOR]` codes:

- `[?GN]`, `[?GREEN]` - green
- `[?BGN]`, `[?BGREEN]` - bright green
- `[?BACKGREEN]`, `[?BKGN]` - green background
- `[?BACKBGREEN]`, `[?BKBGN]` - bright green background
- `[?BOLD]`, `[?BD]` - bold
- `[?DIM]`, `[?DM]` - dim
- `[?ITALIC]`, `[?IT]` - italic
- `[?UNDERLINE]`, `[?UE]` - underlined
- `[?BLINK]`, `[?BL]` - blinking
- `[?REVERSE]`, `[?RV]` - reversed
- `[?STRIKETHROUGH]`, `[?ST]` - strikethrough
- `[?RT]` - reset

Format: `[?<BACKGROUND><BRIGHT><COLOR>]`.

- `<BACKGROUND>` - prefix `BK` or `BACK` for background, optional.
- `<BRIGHT>` - prefix `B` for the bright variant, optional.
- `<COLOR>` - name or two letters: the first and last letters of the name, written uppercase. `GN` = green.

### Palette

**Standard 16**: black, red, green, yellow, blue, magenta, cyan, white and their bright variants. All are available as backgrounds.

**Extended xterm-256** (also available as backgrounds):

| Name     | Abbrev. | Name     | Abbrev. |
|----------|---------|----------|---------|
| orange   | `OE`    | lime     | `LE`    |
| pink     | `PK`    | teal     | `TL`    |
| purple   | `PE`    | navy     | `NY`    |
| violet   | `VT`    | gold     | `GD`    |
| brown    | `BN`    | silver   | `SR`    |
| gray     | `GY`    | maroon   | `MN`    |

Examples: `[?OE]` - orange, `[?BACKGOLD]` / `[?BKGD]` - gold background, `[?BKOE]` - orange background.

### Usage

```rust
use tap::tooling::colors::print::println;
use tap::tooling::colors::set::colorize;

let raw = "[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]";

println(raw);                        // prints already colored
let colored: String = colorize(raw); // or get a string
```

`println` and `print` in `colors::print` already call `colorize` inside - no need to color twice.

Globally disable colors:

```rust
use tap::tooling::colors::set::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

Unknown codes `[?FOO]` are not consumed - they remain in the string as is.

## Grouping commands in help

Commands with identical `CmdInfo` are automatically grouped into one line as aliases:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

Help will show `[help / -h]`.

`required`, `optional`, `unlimited`, and `docs` are compared. If any of them differs, there will be no grouping.

## Renderers

The `tooling::renderable` module templates output by scope:

```rust
use tap::tooling::renderable::{
    Renderable, add_template_render, add_template_render_colored,
};

let mut r = Renderable::new();
r.scope.insert("name", "Rust");
r.scope.insert("version:int", "1");
r.scope.insert("stable:bool", "true");

add_template_render("Hello, {{name}} (v{{version}})", "plain", &mut r);
add_template_render_colored("[?GN]Hello, {{name}}[?RT]", "colored", &mut r);

println!("{}", r.render_with("plain").unwrap());
println!("{}", r.render_with("colored").unwrap());
println!("{}", r.render_json().unwrap());   // {"name":"Rust","stable":true,"version":1}
```

The key suffix sets the type during JSON serialization:

- `:int` - integer; the string must match `-?(0|[1-9][0-9]*)`.
- `:float` - floating-point number with optional exponent.
- `:bool` - only `"true"` or `"false"`.
- `:raw` - the value is inserted into JSON as is, without quotes or validation.
- Without a suffix - string.

`render_json` returns `Result`. On an invalid value it returns `RenderError::InvalidValue`, rather than silently inserting garbage into JSON.

For an array of strings there is `render_json_array(&[Renderable])`:

```rust
use tap::tooling::renderable::render_json_array;

let json = render_json_array(&rows)?;   // [{"...":...},{"...":...}]
```

## Example

The repository has `examples/demo.rs`. It demonstrates all the main features: macros, flags, typed arguments, subcommands, colors, renderers, and JSON.

```bash
cargo run --example demo -- help
cargo run --example demo -- greet Rust
cargo run --example demo -- greet --loud Rust
cargo run --example demo -- add 2 3.5
cargo run --example demo -- serve 127.0.0.1 8080 --verbose --tls=on --workers=8
cargo run --example demo -- report
cargo run --example demo -- echo hello world
```

---

> MIT - see the LICENSE file for details.
> 2026, By Pt.