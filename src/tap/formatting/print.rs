use crate::formatting::set::colorize;


pub fn print(text: &str) {
    print!("{}", colorize(text));
}

pub fn println(text: &str) {
    print!("{}\n", colorize(text));
}