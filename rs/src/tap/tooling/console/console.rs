use crate::tooling::console::{colors::colorize, styles::{Str, Text}};

pub struct Console {
    pub buf: Vec<Text>,
    pub colored: bool,
}

impl Console {

    pub fn new(colored: bool) -> Self {
        Self { buf: vec![], colored: colored }
    }

    pub fn addc(&mut self, text: &str, new: bool) -> &mut Self {
        let mut toappend = Text::new();
        colorize(text, &mut toappend);
        self.addraw(&mut toappend, new)
    }

    pub fn add(&mut self, text: &str, new: bool) -> &mut Self {
        let mut toappend = Text::new();
        toappend.text.push(Str::new().add(text).clone());
        self.addraw(&mut toappend, new)
    }

    fn addraw(&mut self, text: &mut Text, new: bool) -> &mut Self {
        if new || self.buf.is_empty() {
            self.buf.push(text.clone());
        } else {
            let last = self.buf.last_mut().unwrap();
            last.text.append(&mut text.text);
        }
        self
    }

    pub fn clear(&mut self) -> &mut Self {
        self.buf = vec![];
        self
    }

    pub fn get(&self, colored: bool) -> String {
        let mut lines: Vec<String> = vec![];
        for line in self.buf.iter() {
            lines.push(line.get(colored));
        }
        lines.join("\n")
    }

    pub fn print(&self) {
        print!("{}", self.get(self.colored));
    }
}