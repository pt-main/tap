//! ANSI color formatting for terminal output.
//! Delegates to `tooling::console::colors` as the single source of truth.

use std::collections::HashMap;
use std::sync::atomic::AtomicBool;
use std::sync::OnceLock;

use crate::tooling::console::colors;

pub static COLOR_ENABLED: AtomicBool = AtomicBool::new(true);

/// Lookup table: short code (uppercase) -> ANSI escape sequence.
/// Uses `console::colors` as the single source of truth.
pub fn colors() -> &'static HashMap<&'static str, &'static str> {
    static MAP: OnceLock<HashMap<&'static str, &'static str>> = OnceLock::new();
    MAP.get_or_init(|| {
        let mut m = HashMap::with_capacity(128);
        for (code, color) in colors::colors().iter() {
            m.insert(*code, color.code);
        }
        m
    })
}