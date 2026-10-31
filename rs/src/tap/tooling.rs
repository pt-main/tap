pub mod colors;
pub mod console;
pub mod formatting;
pub mod renderable;
use std::rc::Rc;

use crate::{
    engine::structs::{CmdInfo, Engine, ErrorType},
    tooling::formatting::help::generate_help,
};

pub fn help_handler(e: &mut Engine, _: &[&str]) -> ErrorType {
    println!("{}", generate_help(e));
    Ok(())
}

pub fn add_help_handler<'a>(names: Vec<&'a str>, e: &mut Engine<'a>) {
    e.new_commands(
        names,
        CmdInfo::new(
            Rc::new(help_handler),
            vec![],
            vec![],
            false,
            "Generate and print help message",
        ),
    );
}
