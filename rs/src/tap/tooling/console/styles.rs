#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Color {
    pub name: &'static str,
    pub code: &'static str,
}

#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub struct Colors {
    codes: Vec<Color>,
}

impl Colors {
    pub fn new() -> Self {
        Self { codes: Vec::new() }
    }

    pub fn add(&mut self, c: Color) -> &mut Self {
        self.codes.push(c);
        self
    }

    pub fn rem(&mut self, name: &str) -> &mut Self {
        self.codes.retain(|c| c.name != name);
        self
    }

    pub fn get(&self) -> String {
        let mut res = String::new();
        for col in &self.codes {
            res.push_str(col.code);
        }
        res
    }

    pub fn is_empty(&self) -> bool {
        self.codes.is_empty()
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Str {
    pub text: String,
    pub colors: Colors,
}

impl Str {
    pub fn new() -> Self {
        Self { text: "".to_string(), colors: Colors::new() }
    }

    pub fn add(&mut self, text: &str) -> &mut Self {
        self.text += text;
        self
    }

    pub fn get(&self, colored: bool) -> String {
        let mut res = String::from("");
        if colored {
            res += self.colors.get().as_str();
        }
        res += self.text.as_str();
        res
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Text {
    pub text: Vec<Str>,
}

impl Text {
    pub fn new() -> Self {
        Self { text: vec![] }
    }

    pub fn get(&self, colored: bool) -> String {
        let mut res = String::from("");
        for text in self.text.iter() {
            res.push_str(&text.get(colored));
        } 
        res
    }
}