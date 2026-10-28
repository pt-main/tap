use std::rc::Rc;

use crate::engine::structs::{CmdInfo, Engine, ErrorType};

fn validate_args<H>(info: &CmdInfo<'_, H>, args: &[&str]) -> Option<String> {
    let ln = args.len();
    let r = &info.required;
    let o = &info.optional;
    let u = info.unlimited;
    if ln < r.len() {
        return Some(format!(
            "Missing arguments: expected at least {}, got {} ({})",
            r.len(),
            ln,
            r.join(", ")
        ));
    }
    if ln > r.len() + o.len() && !u {
        return Some(format!(
            "Too many arguments: expected at most {}, got {}",
            r.len() + o.len(),
            ln
        ));
    }
    None
}

impl Engine<'_> {
    pub fn process(&mut self, cmd: &str, args: &[&str]) -> ErrorType {
        let err: ErrorType;
        if let Some(info) = self.cmds.get(cmd) {
            if let Some(err) = validate_args(info, args) {
                return Err(err.into());
            }
            let handler = Rc::clone(&info.handler);
            err = handler(self, args);
        } else if let Some(info) = self.subcmds.get(cmd) {
            if let Some(err) = validate_args(info, args) {
                return Err(err.into());
            }
            err = (info.handler)(args);
        } else if let Some(info) = &self.basic_handler {
            let mut _args = vec![cmd];
            _args.extend(args);
            if let Some(err) = validate_args(info, &_args) {
                return Err(err.into());
            }
            let handler = Rc::clone(&info.handler);
            err = handler(self, &_args);
        } else {
            err = Err("Invalid command: not found".into());
        }
        match err {
            Err(e) => return Err(format!("Process Error: {}", e).into()),
            Ok(_) => return Ok(()),
        }
    }
}
