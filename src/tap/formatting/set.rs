use std::sync::atomic::Ordering;

use crate::formatting::color::{COLOR_ENABLED, colors};

pub fn colorize(s: &str) -> String {
    let mut res = s.to_string(); 
    let colors = colors();  
    if COLOR_ENABLED.load(Ordering::Relaxed) {
        for (key, val) in colors.iter() {
            let plhdr = create_placeholder(key); 
            res = res.replace(&plhdr, val);      
        }
    } else {
        for key in colors.keys() {
            let plhdr = create_placeholder(key); 
            res = res.replace(&plhdr, "");      
        }
    }
    res
}

fn create_placeholder(color: &str) -> String {
    return format!("[?{}]", color);
}