use crate::{
    engine::structs::{CmdInfo, Engine}, 
    formatting::{set::colorize, text::{add_space, frame}}
};

#[derive(PartialEq)]
pub struct HelpCmdInfo<'a, H> {
    pub commands: Vec<String>,
    pub info: &'a CmdInfo<'a, H>
}

impl<'a, H> HelpCmdInfo<'a, H> {
    pub fn new(command: String, info: &'a CmdInfo<'a, H>) -> Self {
        Self { commands: vec![command], info: info }
    }

    pub fn is(&mut self, name: String, info: &CmdInfo<'a, H>) -> bool {
        let same = self.info.required == info.required
            && self.info.optional == info.optional
            && self.info.unlimited == info.unlimited
            && self.info.docs == info.docs;
        if same {
            self.commands.push(name);
        }
        same
    }
}


pub fn group_commands<'a, C, H>(commands: &'a C) -> Vec<HelpCmdInfo<'a, H>>
where
    &'a C: IntoIterator<Item = (&'a &'a str, &'a CmdInfo<'a, H>)>,
{
    let mut grouped: Vec<HelpCmdInfo<H>> = Vec::new();
    for (name, info) in commands {
        let mut added = false;
        for cmd in &mut grouped {
            if cmd.is(name.to_string(), info) {
                added = true;
                break;
            }
        }
        if !added {
            grouped.push(HelpCmdInfo::new(name.to_string(), info));
        }
    }
    grouped
}

pub fn generate_help(eng: &Engine) -> String {
    let mut res = String::from("");
    // === Paste about ===
    if eng.about.trim() != "" {
        res.push_str(&format!("{}", colorize(&eng.about)));
        res.push_str("\n\n");
    }
    // === Paste basic handler docs after about ===
    match eng.basic_handler {
        Some(ref i) => {
            let docs = i.docs;
            if docs.trim() != "" {
                res.push_str(docs);
                res.push_str("\n\n");
            }
        }, None => {}
    }
    // === Generating commands help ===
    let commands = group_commands(&eng.cmds);
    for cmd in commands {
        let mut argsv: Vec<String> = vec![];
        for req in &cmd.info.required {
            let r = format!("[?RD]<{}>[?RT]", req.to_string());
            argsv.push(r);
        }
        for opt in &cmd.info.optional {
            let r = format!("[?BE][{}][?RT]", opt.to_string());
            argsv.push(r);
        }
        if cmd.info.unlimited {
            argsv.push("[?CN](arg...)[?RT]".to_string());
        }
        if argsv.len() == 0 {
            argsv.push("(none)".to_string());
        }
        let args = colorize(&argsv.join(", "));
        let block = &format!(
"- [?YW]Args[?RT]: {}
- [?YW]Docs[?RT]:
{}", 
args, add_space(&cmd.info.docs, " ", 1));
        let cmd = &colorize(&format!(
            "[?BGN] Command [?BKBK][?BBK][[?BKBK][?YW]{}[?BBK]][?RT]", 
        cmd.commands.join(" [?GN]OR[?YW] ")));
        res.push_str(&frame(block, cmd, "", "GN"));
        res.push('\n');
    }
    // === Generating subcommands help ===
    let subcommands = group_commands(&eng.subcmds);
    for cmd in subcommands {
        let block = &format!("- [?YW]Docs[?RT]:
{}", add_space(&cmd.info.docs, " ", 1));
        let cmd = &colorize(&format!(
            "[?BGN] Command [?BKBK][?BBK][[?BKBK][?YW]{}[?BBK]][?RT]", 
        cmd.commands.join(" [?GN]OR[?YW] ")));
        res.push_str(&frame(block, cmd, "", "GN"));
        res.push('\n');
    }
    res
}