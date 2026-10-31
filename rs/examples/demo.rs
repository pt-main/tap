//! Tap demo — the main library features in one small CLI.
//!
//! Try:
//!     cargo run --example demo -- help
//!     cargo run --example demo -- version
//!     cargo run --example demo -- greet Rust
//!     cargo run --example demo -- greet --loud Rust
//!     cargo run --example demo -- add 2 3.5
//!     cargo run --example demo -- serve 127.0.0.1 8080 --verbose --tls=on --workers=8
//!     cargo run --example demo -- report
//!     cargo run --example demo -- echo hello world

use tap::{
    argsparsing::ArgsP,
    args, commands, cprintln, ensure, scope, subcommands, template, template_colored,
    engine::structs::{Engine, ErrorType},
    tooling::{add_help_handler, renderable::Renderable},
};

// ── Handlers ─────────────────────────────────────────────────────────

/// `greet [name] [--loud]` — friendly greeting.
fn greet_handler(e: &mut Engine, raw: &[&str]) -> ErrorType {
    args! { raw;
        name: String = "World".to_string(),
    }
    cprintln!("[?BGN]Hello, [?YW]{}[?BGN]![?RT]", name);
    if e.flag("loud") {
        cprintln!("[?BD][?RD]>>> LOUD MODE ON <<<[?RT]");
    }
    Ok(())
}

/// `add <a> <b>` — typed parsing with helpful errors on bad input.
fn add_handler(_: &mut Engine, raw: &[&str]) -> ErrorType {
    args! { raw;
        a: f64,
        b: f64,
    }
    cprintln!("[?GN]{} + {} = [?BGN]{}[?RT]", a, b, a + b);
    Ok(())
}

/// `serve [host] [port] [--verbose] [--tls=on|off|auto] [--workers=N]`
fn serve_handler(e: &mut Engine, raw: &[&str]) -> ErrorType {
    args! { raw;
        host: String = "127.0.0.1".to_string(),
        port: u16 = 8080,
    }
    ensure!(port > 0, "port must be positive (got {})", port);

    let verbose = e.flag("verbose");
    let tls = e.value("tls").unwrap_or("off");
    ensure!(
        matches!(tls, "on" | "off" | "auto"),
        "invalid --tls value {:?} (expected on|off|auto)",
        tls
    );
    let workers: u16 = e.value_t("workers")?.unwrap_or(4);
    ensure!(workers > 0, "--workers must be positive");

    cprintln!("[?YW]Starting server...[?RT]");
    cprintln!("  [?CN]host[?RT]    = {}", host);
    cprintln!("  [?CN]port[?RT]    = {}", port);
    cprintln!("  [?CN]tls[?RT]     = {}", tls);
    cprintln!("  [?CN]workers[?RT] = {}", workers);
    cprintln!("  [?CN]verbose[?RT] = {}", verbose);
    Ok(())
}

/// `report` — render templates and JSON via `tooling::renderable`.
fn report_handler(_: &mut Engine, _: &[&str]) -> ErrorType {
    let mut r = Renderable::new();
    scope!(r;
        "name"        => "Tap",
        "version:int" => "156",
        "stable:bool" => "true",
        "tag"         => "tap-rs",
    );
    template!(r, "plain", "Hello, {{name}} v{{version}} ({{tag}})");
    template_colored!(r, "colored", "[?BGN]Hello, {{name}} [?YW]v{{version}}[?RT]");

    cprintln!("[?DM]--- plain ---[?RT]");
    println!("{}", r.render_with("plain")?);
    cprintln!("[?DM]--- colored ---[?RT]");
    println!("{}", r.render_with("colored")?);
    cprintln!("[?DM]--- json ---[?RT]");
    println!("{}", r.render_json()?);
    Ok(())
}

/// `version` — command with no arguments.
fn version_handler(_: &mut Engine, _: &[&str]) -> ErrorType {
    cprintln!("[?BGN]tap-rs[?RT] v[?YW]1.6.1[?RT]");
    Ok(())
}

/// `echo ...` — subcommand: plain `fn(&[&str])`, no `Engine` access.
fn echo_handler(raw: &[&str]) -> ErrorType {
    ensure!(!raw.is_empty(), "nothing to echo");
    cprintln!("[?GN]echo:[?RT] {}", raw.join(" "));
    Ok(())
}

// ── Wiring ───────────────────────────────────────────────────────────

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "TAP DEMO — lightweight CLI parsing".to_string());

    commands! {
        engine;
        "greet"   => greet_handler   { optional: ["name"],         docs: "Greet someone. Flag: --loud" }
        "add"     => add_handler     { required: ["a", "b"],       docs: "Sum two numbers" }
        "serve"   => serve_handler   { optional: ["host", "port"], docs: "Start a server. Flags: --verbose, --tls=on|off|auto, --workers=N" }
        "report"  => report_handler  {                             docs: "Render templates and JSON" }
        "version" => version_handler {                             docs: "Print version" }
    }

    subcommands! {
        engine;
        "echo" | "say" => echo_handler {
            unlimited: true,
            docs: "Echo back all arguments (subcommand)."
        }
    }

    add_help_handler(vec!["help", "-h"], &mut engine);

    if let Err(e) = engine.main() {
        cprintln!("[?RD][error][?RT] {}", e);
        std::process::exit(1);
    }
}