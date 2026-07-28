use tap::{argsparsing::ArgsP, color, engine::structs::{CmdInfo, Engine}};

fn main() {
    let mut ap = ArgsP::new();
    let mut e = Engine::new(&mut ap);
    e.new_command("cmd", CmdInfo::new(|_, a: &[&str]|{
        println!("args: {:?}", a);
        None
    }, vec![], vec![], true));
    println!("{:?}, {:?}", e.work(&vec!["--test", "cmd", "args", "--val:test", "-val2=val"]), e.args_p.values);
    println!("{}", color::set::set("[s/C+][c/GN]test[c/RT]"))
}
