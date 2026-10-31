use std::rc::Rc;

use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType, HandlerType},
    tooling::{
        colors::print::println as cprintln,
        console::engine::add_help_handler,
        renderable::{Renderable, add_template_render, add_template_render_colored},
    },
};

// ---- handlers ---------------------------------------------------------------

fn greet_handler(e: &mut Engine, args: &[&str]) -> ErrorType {
    let name = args.first().copied().unwrap_or("World");
    cprintln(&format!("[?BGN]Hello, [?YW]{}[?BGN]![?RT]", name));
    if e.args_p.flags.iter().any(|f| f == "loud") {
        cprintln("[?BD][?RD]>>> LOUD MODE ON <<<[?RT]");
    }
    Ok(())
}

fn add_handler(_: &mut Engine, args: &[&str]) -> ErrorType {
    let a: f64 = args[0]
        .parse()
        .map_err(|_| format!("bad number: {}", args[0]))?;
    let b: f64 = args[1]
        .parse()
        .map_err(|_| format!("bad number: {}", args[1]))?;
    cprintln(&format!("[?GN]{} + {} = [?BGN]{}[?RT]", a, b, a + b));
    Ok(())
}

fn serve_handler(e: &mut Engine, args: &[&str]) -> ErrorType {
    let host = args[0];
    let port = args[1];
    let verbose = e.args_p.flags.iter().any(|f| f == "verbose");
    let tls = e
        .args_p
        .values
        .get("tls")
        .map(String::as_str)
        .unwrap_or("off");
    cprintln("[?YW]Starting server...[?RT]");
    cprintln(&format!("  [?CN]host[?RT]    = {}", host));
    cprintln(&format!("  [?CN]port[?RT]    = {}", port));
    cprintln(&format!("  [?CN]tls[?RT]     = {}", tls));
    cprintln(&format!("  [?CN]verbose[?RT] = {}", verbose));
    Ok(())
}

fn report_handler(_: &mut Engine, _: &[&str]) -> ErrorType {
    let mut r = Renderable::new();
    r.scope.insert("name", "Tap");
    r.scope.insert("version:int", "156");
    r.scope.insert("stable:bool", "true");
    r.scope.insert("tag", "tap-rs");

    add_template_render("Hello, {{name}} v{{version}}", "plain", &mut r);
    add_template_render_colored(
        "[?BGN]Hello, {{name}} [?YW]v{{version}}[?RT]",
        "colored",
        &mut r,
    );

    cprintln("[?DM]--- plain ---[?RT]");
    println!("{}", r.render_with("plain")?);
    cprintln("[?DM]--- colored ---[?RT]");
    println!("{}", r.render_with("colored")?);
    cprintln("[?DM]--- json ---[?RT]");
    println!("{}", r.render_json()?);
    Ok(())
}

// ---- wiring -----------------------------------------------------------------

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "TAP DEMO — simple but powerful".to_string());

    let greet_h: HandlerType = Rc::new(greet_handler);
    let add_h: HandlerType = Rc::new(add_handler);
    let serve_h: HandlerType = Rc::new(serve_handler);
    let report_h: HandlerType = Rc::new(report_handler);

    engine.new_command(
        "greet",
        CmdInfo::new(
            greet_h,
            vec![],
            vec!["name"],
            false,
            "Greet someone. Optional: name. Flag: --loud",
        ),
    );
    engine.new_command(
        "add",
        CmdInfo::new(add_h, vec!["a", "b"], vec![], false, "Sum two numbers"),
    );
    engine.new_command(
        "serve",
        CmdInfo::new(
            serve_h,
            vec!["host", "port"],
            vec![],
            false,
            "Start a server. Flags: --verbose, --tls=on|off",
        ),
    );
    engine.new_command(
        "report",
        CmdInfo::new(report_h, vec![], vec![], false, "Render templates + JSON"),
    );

    add_help_handler(vec!["help", "-h", "--help"], &mut engine);

    if let Err(e) = engine.main() {
        eprintln!("[error] {}", e);
        std::process::exit(1);
    }
}
