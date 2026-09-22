use crate::argsparsing::ArgsP;
use std::collections::HashMap;

pub type ErrorType = Result<(), Box<dyn std::error::Error>>;
pub type HandlerType = fn(&mut Engine<'_>, &[&str]) -> ErrorType;
pub type SubHandlerType = fn(&[&str]) -> ErrorType;

pub type CmdInfoCmd<'a> = CmdInfo<'a, HandlerType>;
pub type CmdInfoSubcmd<'a> = CmdInfo<'a, SubHandlerType>;

#[derive(PartialEq, Clone)]
pub struct CmdInfo<'a, H> {
    pub handler: H,
    pub required: Vec<&'a str>,
    pub optional: Vec<&'a str>,
    pub unlimited: bool,
    pub docs: &'static str,
}

pub type CmdMap<'a> = HashMap<&'a str, CmdInfoCmd<'a>>;
pub type SubCmdsMap<'a> = HashMap<&'a str, CmdInfoSubcmd<'a>>;
pub type BasicHandler<'a> = Option<CmdInfo<'a, HandlerType>>;

pub struct Engine<'a> {
    pub args_p: &'a mut ArgsP,
    pub cmds: CmdMap<'a>,
    pub subcmds: SubCmdsMap<'a>,
    pub basic_handler: BasicHandler<'a>,
    pub about: String,
}