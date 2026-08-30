use crate::argsparsing::ArgsP;
use crate::engine::structs::{BasicHandler, CmdInfoCmd, CmdInfoSubcmd, CmdMap, Engine, ErrorType, SubCmdsMap};
use std::collections::HashMap;
use std::{env};

pub struct EngineGetter<'a> {
    pub e: &'a Engine<'a>
}

impl<'a> EngineGetter<'a> {
    pub fn new(e: &'a Engine<'a>) -> Self {
        Self { e: e }
    }

    pub fn get_cmds(&self) -> &CmdMap<'a> {
        return &self.e.cmds;
    }

    pub fn get_subcmds(&self) -> &SubCmdsMap<'a> {
        return &self.e.subcmds;
    }

    pub fn basic_handler(&self) -> &BasicHandler<'a> {
        return  &self.e.basic_handler;
    }
}

impl<'a> Engine<'a> {
    pub fn new(args_p: &'a mut ArgsP, about: String) -> Self {
        Self {
            args_p,
            cmds: HashMap::new(),
            subcmds: HashMap::new(),
            basic_handler: None,
            about: about,
        }
    }

    pub fn getter(&'a self) -> EngineGetter<'a> {
        return EngineGetter::new(self);
    }

    pub fn new_command(&mut self, name: &'a str, cmdinfo: CmdInfoCmd<'a>) {
        self.cmds.insert(name, cmdinfo);
    }

    pub fn new_commands(&mut self, names: Vec<&'a str>, cmdinfo: CmdInfoCmd<'a>) {
        for name in names {
            self.new_command(name, cmdinfo.clone());
        }
    }

    pub fn new_subcommand(&mut self, name: &'a str, cmdinfo: CmdInfoSubcmd<'a>) {
        self.subcmds.insert(name, cmdinfo);
    }

    pub fn new_subcommands(&mut self, names: Vec<&'a str>, cmdinfo: CmdInfoSubcmd<'a>) {
        for name in names {
            self.new_subcommand(name, cmdinfo.clone());
        }
    }

    pub fn work(&mut self, args: &[&str]) -> ErrorType {
        self.args_p.parse_agrs_to_self(args);
        if self.args_p.args.is_empty() {
            println!("{}", self.about);
            return Err("No command provided".into());
        }
        let all_args = self.args_p.args.clone();
        let cmd = all_args[0].as_str();
        let args_refs: Vec<&str> = all_args[1..].iter().map(|s| s.as_str()).collect();
        self.process(cmd, &args_refs)
    }

    pub fn main(&mut self) -> ErrorType {
        self.work(
            &env::args()
                .collect::<Vec<_>>()
                .iter()
                .map(|s| s.as_str())
                .collect::<Vec<_>>()[1..],
        )
    }
}
