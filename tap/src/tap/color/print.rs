use crate::color::set::set;

pub fn print(text: &str) {
    print!("{}", set(text));
}

pub fn println(text: &str) {
    print(&(text.to_owned() + "\n"))
}