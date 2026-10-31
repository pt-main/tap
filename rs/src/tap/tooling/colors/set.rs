use std::sync::atomic::{AtomicBool, Ordering};

use crate::tooling::console::{colors, styles::Text};

pub static COLOR_ENABLED: AtomicBool = AtomicBool::new(true);

pub fn colorize(s: &str) -> String {
    let colored = COLOR_ENABLED.load(Ordering::Relaxed);
    let mut text = Text::new();
    colors::colorize(s, &mut text);
    text.get(colored)
}
