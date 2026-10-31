//! Convenience macros and helper functions for `tap`.
//!
//! # Enabling
//!
//! This module is not wired up automatically. Add it in `src/tap/mod.rs`:
//!
//! ```ignore
//! pub mod macros;
//! ```
//!
//! All macros are exported at the crate root (`#[macro_export]`), so call them
//! through the crate name — e.g. `tap::commands! { ... }`, `tap::bail!(...)`,
//! `tap::cprintln!(...)`.
//!
//! # Contents
//!
//! | Category             | Items                                                                    |
//! |----------------------|--------------------------------------------------------------------------|
//! | Helpers              | [`arg`], [`arg_str`], [`arg_or`], [`arg_str_or`], [`arg_str_opt`]        |
//! | Registration         | [`commands!`], [`subcommands!`]                                          |
//! | Argument parsing     | [`args!`]                                                                |
//! | Errors               | [`bail!`], [`ensure!`]                                                   |
//! | Output               | [`cprint!`], [`cprintln!`], [`ceprintln!`]                               |
//! | Scopes / templates   | [`scope!`], [`template!`], [`template_colored!`]                         |
//!
//! # Typical usage
//!
//! ```ignore
//! use tap::{commands, args, bail, ensure};
//!
//! commands! {
//!     engine;
//!     "add" | "a" => my_add {
//!         required: ["x", "y"],
//!         docs: "Add two numbers",
//!     }
//!     "greet" => my_greet {
//!         required: ["name"],
//!         optional: ["greeting"],
//!         docs: "Say hello",
//!     }
//! }
//!
//! fn my_add(_: &mut Engine, args: &[&str]) -> ErrorType {
//!     args! { args;
//!         x: i64,
//!         y: i64,
//!     }
//!     println!("{}", x + y);
//!     Ok(())
//! }
//! ```

use std::str::FromStr;

// ─────────────────────────────────────────────────────────────────────
//  Helper functions
// ─────────────────────────────────────────────────────────────────────

/// Parses a required argument at index `i` into type `T`.
///
/// Returns an error when:
///
/// * no argument exists at index `i` — `"missing argument <name> (position i)"`;
/// * the argument exists but fails to parse as `T` — `"argument <name>: <parser error>"`.
///
/// The error message includes `name` so users can tell which argument was
/// wrong rather than receiving a bare parser error.
pub fn arg<'a, T>(args: &[&'a str], i: usize, name: &str) -> Result<T, String>
where
    T: FromStr,
    T::Err: std::fmt::Display,
{
    let s = args
        .get(i)
        .ok_or_else(|| format!("missing argument <{}> (position {})", name, i))?;
    s.parse()
        .map_err(|e| format!("argument <{}>: {}", name, e))
}

/// Returns a required argument at index `i` as a string slice.
///
/// Unlike [`arg`], no parsing is performed. Errors only when the argument is
/// missing.
pub fn arg_str<'a>(args: &[&'a str], i: usize, name: &str) -> Result<&'a str, String> {
    args.get(i)
        .copied()
        .ok_or_else(|| format!("missing argument <{}> (position {})", name, i))
}

/// Like [`arg`], but falls back to `default` when the argument is absent.
///
/// Parsing errors are still reported — if the argument is present, it must
/// parse as `T`.
pub fn arg_or<'a, T>(args: &[&'a str], i: usize, name: &str, default: T) -> Result<T, String>
where
    T: FromStr,
    T::Err: std::fmt::Display,
{
    match args.get(i) {
        Some(s) => s.parse().map_err(|e| format!("argument <{}>: {}", name, e)),
        None => Ok(default),
    }
}

/// Like [`arg_str`], but returns `default` when the argument is absent.
///
/// Never fails.
pub fn arg_str_or<'a>(args: &[&'a str], i: usize, default: &'a str) -> &'a str {
    args.get(i).copied().unwrap_or(default)
}

/// Returns the argument at index `i` if present, otherwise `None`.
///
/// Never fails.
pub fn arg_str_opt<'a>(args: &[&'a str], i: usize) -> Option<&'a str> {
    args.get(i).copied()
}

// ─────────────────────────────────────────────────────────────────────
//  commands!
// ─────────────────────────────────────────────────────────────────────

/// Registers one or more top-level commands on an engine.
///
/// Each entry has the shape:
///
/// ```ignore
/// "name" | "alias" => handler {
///     required: ["arg1", "arg2"],   // optional field
///     optional: ["arg3"],           // optional field
///     unlimited: true,              // optional, defaults to false
///     docs: "Description",          // required
/// }
/// ```
///
/// All fields except `docs` may be omitted. `required`, `optional` and
/// `unlimited` default to empty/false.
///
/// # Example
///
/// ```ignore
/// tap::commands! {
///     engine;
///     "echo" => cmd_echo {
///         required: ["text"],
///         docs: "Echo back the text",
///     }
///     "quit" | "exit" => cmd_quit {
///         docs: "Exit the shell",
///     }
/// }
/// ```
#[macro_export]
macro_rules! commands {
    (
        $engine:ident;
        $(
            $($name:literal)|+ => $handler:path {
                $(required: [$($req:literal),* $(,)?],)?
                $(optional: [$($opt:literal),* $(,)?],)?
                $(unlimited: $unlim:literal,)?
                docs: $docs:literal $(,)?
            }
        )*
    ) => {{
        $(
            $engine.new_commands(
                vec![$($name),+],
                {
                    let __tap_h: $crate::engine::structs::HandlerType =
                        ::std::rc::Rc::new($handler);
                    $crate::engine::structs::CmdInfo::new(
                        __tap_h,
                        vec![$($($req),*)?],
                        vec![$($($opt),*)?],
                        $crate::commands!(@unlim $($unlim)?),
                        $docs,
                    )
                },
            );
        )*
    }};

    (@unlim $b:literal) => { $b };
    (@unlim)            => { false };
}

// ─────────────────────────────────────────────────────────────────────
//  subcommands!
// ─────────────────────────────────────────────────────────────────────

/// Registers one or more subcommands on an engine.
///
/// Syntax is identical to [`commands!`], except that the handler must be a
/// `SubHandlerType` (a plain function pointer, not an `Rc`).
///
/// # Example
///
/// ```ignore
/// tap::subcommands! {
///     engine;
///     "set" => sub_set {
///         required: ["key", "value"],
///         docs: "Set a variable",
///     }
/// }
/// ```
#[macro_export]
macro_rules! subcommands {
    (
        $engine:ident;
        $(
            $($name:literal)|+ => $handler:path {
                $(required: [$($req:literal),* $(,)?],)?
                $(optional: [$($opt:literal),* $(,)?],)?
                $(unlimited: $unlim:literal,)?
                docs: $docs:literal $(,)?
            }
        )*
    ) => {{
        $(
            $engine.new_subcommands(
                vec![$($name),+],
                $crate::engine::structs::CmdInfo::new(
                    $handler as $crate::engine::structs::SubHandlerType,
                    vec![$($($req),*)?],
                    vec![$($($opt),*)?],
                    $crate::subcommands!(@unlim $($unlim)?),
                    $docs,
                ),
            );
        )*
    }};

    (@unlim $b:literal) => { $b };
    (@unlim)            => { false };
}

// ─────────────────────────────────────────────────────────────────────
//  args!
// ─────────────────────────────────────────────────────────────────────

/// Binds handler arguments to typed local variables.
///
/// Syntax: `args! { <slice>; name: Type, name: Type = default, ... }`.
///
/// * `name: Type`             — required, parsed via [`arg`].
/// * `name: Type = default`   — optional, parsed via [`arg_or`].
///
/// Each variable consumes one positional argument, so the order matters.
/// The trailing comma is allowed and optional.
///
/// # Example
///
/// ```ignore
/// fn handler(_: &mut Engine, args: &[&str]) -> ErrorType {
///     args! { args;
///         count: u32,
///         label: String = "default".to_string(),
///         verbose: bool = false,
///     }
///     // ...
///     Ok(())
/// }
/// ```
#[macro_export]
macro_rules! args {
    ($args:expr; $($rest:tt)*) => {
        let __tap_args: &[&str] = $args;
        let mut __tap_i: usize = 0usize;
        $crate::args!(@next __tap_args, __tap_i, $($rest)*);
        let _ = __tap_i;
    };

    (@next $args:ident, $i:ident,) => {};

    (@next $args:ident, $i:ident, $name:ident : $ty:ty = $default:expr, $($rest:tt)*) => {
        let $name: $ty = $crate::macros::arg_or(
            $args, $i, stringify!($name), $default
        )?;
        $i += 1;
        $crate::args!(@next $args, $i, $($rest)*);
    };

    (@next $args:ident, $i:ident, $name:ident : $ty:ty = $default:expr) => {
        let $name: $ty = $crate::macros::arg_or(
            $args, $i, stringify!($name), $default
        )?;
        $i += 1;
        $crate::args!(@next $args, $i,);
    };

    (@next $args:ident, $i:ident, $name:ident : $ty:ty, $($rest:tt)*) => {
        let $name: $ty = $crate::macros::arg(
            $args, $i, stringify!($name)
        )?;
        $i += 1;
        $crate::args!(@next $args, $i, $($rest)*);
    };

    (@next $args:ident, $i:ident, $name:ident : $ty:ty) => {
        let $name: $ty = $crate::macros::arg(
            $args, $i, stringify!($name)
        )?;
        $i += 1;
        $crate::args!(@next $args, $i,);
    };
}

// ─────────────────────────────────────────────────────────────────────
//  Error helpers — `bail!` / `ensure!`
// ─────────────────────────────────────────────────────────────────────

/// Early-returns an error from a handler.
///
/// Accepts either a literal/expression convertible into the handler's error
/// type, or a `format!`-style format string with arguments.
///
/// ```ignore
/// fn h(_: &mut Engine, args: &[&str]) -> ErrorType {
///     if args.is_empty() { bail!("at least one argument required"); }
///     let n: u32 = args[0].parse().map_err(|_| format!("bad number: {}", args[0]))?;
///     ensure!(n > 0, "must be positive, got {}", n);
///     Ok(())
/// }
/// ```
#[macro_export]
macro_rules! bail {
    ($msg:literal $(,)?) => {
        return ::std::result::Result::Err($msg.into())
    };
    ($err:expr $(,)?) => {
        return ::std::result::Result::Err($err.into())
    };
    ($fmt:expr, $($arg:tt)*) => {
        return ::std::result::Result::Err(format!($fmt, $($arg)*).into())
    };
}

/// Conditional [`bail!`]: returns an error when the condition is false.
///
/// ```ignore
/// ensure!(count <= 10, "too many items: {}", count);
/// ensure!(flag, "flag must be set");
/// ```
#[macro_export]
macro_rules! ensure {
    ($cond:expr, $msg:literal $(,)?) => {
        if !$cond { $crate::bail!($msg) }
    };
    ($cond:expr, $err:expr $(,)?) => {
        if !$cond { $crate::bail!($err) }
    };
    ($cond:expr, $fmt:expr, $($arg:tt)*) => {
        if !$cond { $crate::bail!($fmt, $($arg)*) }
    };
}

// ─────────────────────────────────────────────────────────────────────
//  Colored printing — `cprint!` / `cprintln!` / `ceprintln!`
// ─────────────────────────────────────────────────────────────────────

/// Colored `print!`.
///
/// The formatted string is passed through
/// `tooling::colors::set::colorize` before writing to stdout.
#[macro_export]
macro_rules! cprint {
    ($($arg:tt)*) => {
        ::std::print!("{}", $crate::tooling::colors::set::colorize(&format!($($arg)*)))
    };
}

/// Colored `println!`.
///
/// With no arguments, emits a plain newline (no colorization needed).
/// Otherwise formats and colorizes the output.
#[macro_export]
macro_rules! cprintln {
    () => { ::std::println!() };
    ($($arg:tt)*) => {
        ::std::println!("{}", $crate::tooling::colors::set::colorize(&format!($($arg)*)))
    };
}

/// Colored `eprintln!` — writes to stderr.
#[macro_export]
macro_rules! ceprintln {
    ($($arg:tt)*) => {
        ::std::eprintln!("{}", $crate::tooling::colors::set::colorize(&format!($($arg)*)))
    };
}

// ─────────────────────────────────────────────────────────────────────
//  Renderable shortcuts — `scope!` / `template!` / `template_colored!`
// ─────────────────────────────────────────────────────────────────────

/// Populates `Renderable::scope` with `key => value` pairs.
///
/// The `key` may include a `:type` suffix (e.g. `"version:int"`) that is
/// interpreted by the renderer downstream.
///
/// ```ignore
/// scope!(r;
///     "name" => "Rust",
///     "version:int" => "1",
///     "stable:bool" => "true",
/// );
/// ```
#[macro_export]
macro_rules! scope {
    ($r:expr; $($key:literal => $val:expr),* $(,)?) => {{
        $($r.scope.insert($key, $val);)*
    }};
}

/// Adds a plain template renderer.
///
/// ```ignore
/// template!(r, "plain", "Hello, {{name}}");
/// ```
#[macro_export]
macro_rules! template {
    ($r:expr, $name:expr, $templ:expr) => {
        $crate::tooling::renderable::add_template_render($templ, $name, &mut $r)
    };
}

/// Adds a colored template renderer.
///
/// The rendered result is passed through `colorize`, so `[?GN]…[?RT]` style
/// color codes are resolved.
///
/// ```ignore
/// template_colored!(r, "colored", "[?GN]Hello, {{name}}[?RT]");
/// ```
#[macro_export]
macro_rules! template_colored {
    ($r:expr, $name:expr, $templ:expr) => {
        $crate::tooling::renderable::add_template_render_colored($templ, $name, &mut $r)
    };
}