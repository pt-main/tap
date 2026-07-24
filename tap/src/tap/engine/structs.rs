use crate::argsparsing::ArgsP;
use std::collections::HashMap;

pub type HandlerType = fn(&mut Engine<'_>, &[&str]) -> Option<String>;
pub type SubHandlerType = fn(&[&str]) -> Option<String>;

pub type CmdInfoCmd<'a> = CmdInfo<'a, HandlerType>;
pub type CmdInfoSubcmd<'a> = CmdInfo<'a, SubHandlerType>;

pub struct CmdInfo<'a, H> {
    pub handler: H,
    pub required: Vec<&'a str>,
    pub optional: Vec<&'a str>,
    pub unlimited: bool,
}

pub struct Engine<'a> {
    pub args_p: &'a mut ArgsP,
    pub cmds: HashMap<&'a str, CmdInfoCmd<'a>>,
    pub subcmds: HashMap<&'a str, CmdInfoSubcmd<'a>>,
    pub basic_handler: Option<CmdInfo<'a, HandlerType>>,
}