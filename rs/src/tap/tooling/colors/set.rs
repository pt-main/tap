use std::sync::atomic::Ordering;

use crate::{tooling::colors::color::{COLOR_ENABLED, colors}, tooling::console::console::Console};

pub fn colorize(s: &str) -> String {
    let mut res = s.to_string(); 
    if COLOR_ENABLED.load(Ordering::Relaxed) {
        res = Console::new(true).addc(s, false).get(true);
    } else {
        for key in colors().keys() {
            let plhdr = create_placeholder(key); 
            res = res.replace(&plhdr, "");      
        }
    }
    res
}

fn create_placeholder(color: &str) -> String {
    return format!("[?{}]", color);
}