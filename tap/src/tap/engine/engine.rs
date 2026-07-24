use crate::argsparsing::ArgsP;
use crate::engine::structs::{Engine, CmdInfoSubcmd, CmdInfoCmd};
use std::collections::HashMap;


impl<'a> Engine<'a> {
    pub fn new(args_p: &'a mut ArgsP) -> Self {
        Self {
            args_p,
            cmds: HashMap::new(),
            subcmds: HashMap::new(),
            basic_handler: None,
        }
    }

    pub fn new_command(&mut self, name: &'a str, cmdinfo: CmdInfoCmd<'a>) {
        self.cmds.insert(name, cmdinfo);
    }

    pub fn new_subcommand(&mut self, name: &'a str, cmdinfo: CmdInfoSubcmd<'a>) {
        self.subcmds.insert(name, cmdinfo);
    }

    pub fn work(&mut self, args: &[&str]) -> Result<(), String> {
        self.args_p.parse_ags_to_self(args);
        if self.args_p.args.is_empty() {
            return Err("No command provided".to_string());
        }
        let all_args = self.args_p.args.clone();
        let cmd = all_args[0].as_str();
        let args_refs: Vec<&str> = all_args[1..].iter().map(|s| s.as_str()).collect();
        match self.process(cmd, &args_refs) {
            Some(val) => Err(val),
            None => Ok(()),
        }
    }
}