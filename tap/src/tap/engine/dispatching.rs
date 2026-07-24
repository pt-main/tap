use crate::engine::structs::{Engine, CmdInfo};


fn validate_args<H>(info: &CmdInfo<'_, H>, args: &[&str]) -> Option<String> {
    let ln = args.len();
    let r = &info.required;
    let o = &info.optional;
    let u = info.unlimited;
    if (ln > r.len() + o.len() && !u) || ln < r.len() {
        return Some(format!(
            "Invalid argument length: {} (must be: [{}])",
            ln, r.join(", ")
        ));
    }
    None
}


impl Engine<'_> {
    pub fn process(&mut self, cmd: &str, args: &[&str]) -> Option<String> {
        let err: Option<String>;
        if let Some(info) = self.cmds.get(cmd) {
            if let Some(err) = validate_args(info, args) {
                return Some(err);
            }
            err = (info.handler)(self, args);
        } else if let Some(info) = self.subcmds.get(cmd) {
            if let Some(err) = validate_args(info, args) {
                return Some(err);
            }
            err = (info.handler)(args);
        } else if let Some(info) = &self.basic_handler {
            if let Some(err) = validate_args(info, args) {
                return Some(err);
            }
            err = (info.handler)(self, args);
        } else {
            err = Some("Invalid command: not found".to_string());
        }
        match err {
            Some(s) => return Some(format!("Process Error: {}", s)),
            None => return None,
        }
    }
}