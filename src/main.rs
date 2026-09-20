use tap::{tooling::colors::set::colorize, tooling::console::console::Console};

fn main() {
    let mut c = Console::new(true);
    c.addc("test ", false);
    c.addc("[?RED]red ", false); 
    c.addc("[?BOLD]bold ", false);
    c.addc("[?UE]text", false);
    c.addc("[?RT]", true);
    c.print();

    println!("{}", colorize("[?BE]TEST[?RT]"))
}