# Tap – a CLI parsing library for Rust

[![Crates.io](https://img.shields.io/crates/v/tap.svg)](https://crates.io/crates/tap)
[![Docs.rs](https://docs.rs/tap/badge.svg)](https://docs.rs/tap)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://opensource.org/licenses/Apache-2.0)

```toml
[dependencies]
tap = "0.1.0"
```

**Tap** is a lightweight library for building CLIs in Rust. It provides a simple API for argument parsing, commands, flags, and colored output, with no external dependencies or complex macros.

---

## Features

- **Colored output** via short codes: `[?GN]`, `[?RD]`, `[?YW]`, etc.
- **Auto‑generated help** with commands grouped by identical handlers.
- **Subcommands** (nested commands).
- **Zero dependencies** (only `std`).

---

## Quick Start

Let's create a simple CLI with a `hello` command:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    utils::add_help_handler,
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

---

## Commands and Arguments

Add a command using `Engine::new_command`:

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> Result<(), Box<dyn Error>> { ... }

engine.new_command(
    "copy",                               // command name
    CmdInfo::new(
        my_handler,                       // handler function
        vec!["src", "dst"],               // required arguments
        vec!["force"],                    // optional arguments
        false,                            // unlimited (if true, accepts any number of extra arguments)
        "Copy source to destination",     // help description
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

---

## Flags

Flags are written as `--flag`, `--key=value`, or `--key:value`. They are parsed automatically and stored in `ArgsP`:

- `args_p.flags` – list of flags without values.
- `args_p.values` – dictionary of key → value.

In your handler you can read them:

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

---

## Colors

Tap supports colored output via short codes like `[?COLOR]`. Examples:
- `[?GN]`, `[?GREEN]` – green
- `[?BGN]`, `[?BGREEN]` – bright green
- `[?BACKGREEN]`, `[?BKGN]` – green background
- `[?BACKBGREEN]`, `[?BKBGN]` – bright green background
- `[?BOLD]`, `[?BD]` – bold
- `[?UNDERLINE]`, `[?UE]` – underlined
- `[?RT]` – reset

Format – `[?<BACKGROUND><BRIGHT><COLOR>]`, where:

- `<BACKGROUND>` – to choose a background colour, prepend `BK` (short) or `BACK` (full name) to the colour.
- `<BRIGHT>` – to use the bright version, prepend `B` (short for 'bright').
- `<COLOR>` – options:
    - Two letters: first letter – initial of the colour name, second – final letter.
        Example: `GN`

Colours – black, red, green, yellow, blue, magenta, cyan, and their bright variants. All colours are available for the background. For text, bold and underline styles are available.

Use functions from `tap::formatting::print` and `tap::formatting::color`:

```rust
use tap::formatting::print::println;
use tap::formatting::color::colorize;

let raw = "[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]";

println(raw); 
let colored = colorize(raw)
```

Colours can be disabled globally:

```rust
use tap::formatting::color::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

---

## Grouping Commands in Help

Commands with identical `CmdInfo` values are automatically grouped into one line as aliases:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

Help will display: `[help / -h]`.

---

> MIT – see the LICENSE file for details.
> 2026, By Pt.