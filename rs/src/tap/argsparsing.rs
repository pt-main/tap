use std::{collections::HashMap, str::FromStr};

pub struct ArgsP {
    pub args: Vec<String>,
    pub flags: Vec<String>,
    pub values: HashMap<String, String>,
}

impl ArgsP {
    pub fn new() -> Self {
        Self {
            args: vec![],
            flags: vec![],
            values: HashMap::new(),
        }
    }

    fn count_leading_same_chars(s: &str) -> usize {
        let mut chars = s.chars();
        match chars.next() {
            Some(first) => 1 + chars.take_while(|&c| c == first).count(),
            None => 0,
        }
    }

    fn parse_args(args: &[&str]) -> (Vec<String>, HashMap<String, String>, Vec<String>) {
        let mut flags: Vec<String> = vec![];
        let mut values: HashMap<String, String> = HashMap::new();
        let mut r_args: Vec<String> = vec![];
        let mut can_be_flag = true;
        for arg in args {
            let mut split: Vec<&str> = arg.splitn(2, "=").collect();
            if split.len() == 1 {
                split = arg.splitn(2, ":").collect();
            }
            let s0 = split[0];
            let sch = Self::count_leading_same_chars(s0);

            if *arg == "--" && can_be_flag {
                can_be_flag = !can_be_flag;
                continue;
            }

            if s0.starts_with("--") && s0.len() > sch && sch < 3 && can_be_flag {
                if split.len() > 1 {
                    values.insert(s0[sch..].to_string(), split[1].to_string());
                } else {
                    flags.push(s0[sch..].to_string());
                }
            } else {
                r_args.push(arg.to_string());
            }
        }
        (flags, values, r_args)
    }

    pub fn parse_args_to_self(&mut self, args: &[&str]) {
        let (flags, values, r_args) = Self::parse_args(args);
        self.flags = flags;
        self.values = values;
        self.args = r_args;
    }
}


pub trait ArgsExt {
    fn arg<T: FromStr>(&self, i: usize) -> Result<T, String>
    where
        T::Err: std::fmt::Display;

    fn arg_or<'a>(&'a self, i: usize, default: &'a str) -> &'a str;
}

impl ArgsExt for [&str] {
    fn arg<T: FromStr>(&self, i: usize) -> Result<T, String>
    where
        T::Err: std::fmt::Display,
    {
        self.get(i)
            .ok_or_else(|| format!("missing argument #{}", i))?
            .parse()
            .map_err(|e| format!("invalid argument #{}: {}", i, e))
    }

    fn arg_or<'a>(&'a self, i: usize, default: &'a str) -> &'a str {
        self.get(i).copied().unwrap_or(default)
    }
}