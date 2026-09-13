//! ANSI color formatting for terminal output.
//! Supports short codes like [?GN], [?RD], [?BOLD], etc.
//! Coloring can be globally disabled by setting COLOR_ENABLED = false.

use std::collections::HashMap;
use std::sync::atomic::AtomicBool;
use std::sync::OnceLock;

pub static COLOR_ENABLED: AtomicBool = AtomicBool::new(true);

// Text foreground
pub const BLACK: &str = "\x1b[30m";
pub const RED: &str = "\x1b[31m";
pub const GREEN: &str = "\x1b[32m";
pub const YELLOW: &str = "\x1b[33m";
pub const BLUE: &str = "\x1b[34m";
pub const MAGENTA: &str = "\x1b[35m";
pub const CYAN: &str = "\x1b[36m";
pub const WHITE: &str = "\x1b[37m";

// Bright foreground
pub const BRIGHT_BLACK: &str = "\x1b[90m";
pub const BRIGHT_RED: &str = "\x1b[91m";
pub const BRIGHT_GREEN: &str = "\x1b[92m";
pub const BRIGHT_YELLOW: &str = "\x1b[93m";
pub const BRIGHT_BLUE: &str = "\x1b[94m";
pub const BRIGHT_MAGENTA: &str = "\x1b[95m";
pub const BRIGHT_CYAN: &str = "\x1b[96m";
pub const BRIGHT_WHITE: &str = "\x1b[97m";

// Extended foreground (xterm-256 approximations)
pub const ORANGE: &str = "\x1b[38;5;208m";
pub const PINK: &str = "\x1b[38;5;213m";
pub const PURPLE: &str = "\x1b[38;5;93m";
pub const VIOLET: &str = "\x1b[38;5;177m";
pub const BROWN: &str = "\x1b[38;5;130m";
pub const GRAY: &str = "\x1b[38;5;245m";
pub const LIME: &str = "\x1b[38;5;118m";
pub const TEAL: &str = "\x1b[38;5;30m";
pub const NAVY: &str = "\x1b[38;5;17m";
pub const GOLD: &str = "\x1b[38;5;220m";
pub const SILVER: &str = "\x1b[38;5;250m";
pub const MAROON: &str = "\x1b[38;5;88m";

// Background (standard)
pub const BLACK_BG: &str = "\x1b[40m";
pub const RED_BG: &str = "\x1b[41m";
pub const GREEN_BG: &str = "\x1b[42m";
pub const YELLOW_BG: &str = "\x1b[43m";
pub const BLUE_BG: &str = "\x1b[44m";
pub const MAGENTA_BG: &str = "\x1b[45m";
pub const CYAN_BG: &str = "\x1b[46m";
pub const WHITE_BG: &str = "\x1b[47m";

// Background (bright)
pub const BRIGHT_BLACK_BG: &str = "\x1b[100m";
pub const BRIGHT_RED_BG: &str = "\x1b[101m";
pub const BRIGHT_GREEN_BG: &str = "\x1b[102m";
pub const BRIGHT_YELLOW_BG: &str = "\x1b[103m";
pub const BRIGHT_BLUE_BG: &str = "\x1b[104m";
pub const BRIGHT_MAGENTA_BG: &str = "\x1b[105m";
pub const BRIGHT_CYAN_BG: &str = "\x1b[106m";
pub const BRIGHT_WHITE_BG: &str = "\x1b[107m";

// Extended background (xterm-256 approximations)
pub const ORANGE_BG: &str = "\x1b[48;5;208m";
pub const PINK_BG: &str = "\x1b[48;5;213m";
pub const PURPLE_BG: &str = "\x1b[48;5;93m";
pub const VIOLET_BG: &str = "\x1b[48;5;177m";
pub const BROWN_BG: &str = "\x1b[48;5;130m";
pub const GRAY_BG: &str = "\x1b[48;5;245m";
pub const LIME_BG: &str = "\x1b[48;5;118m";
pub const TEAL_BG: &str = "\x1b[48;5;30m";
pub const NAVY_BG: &str = "\x1b[48;5;17m";
pub const GOLD_BG: &str = "\x1b[48;5;220m";
pub const SILVER_BG: &str = "\x1b[48;5;250m";
pub const MAROON_BG: &str = "\x1b[48;5;88m";

// Styles
pub const RESET: &str = "\x1b[0m";
pub const BOLD: &str = "\x1b[1m";
pub const DIM: &str = "\x1b[2m";
pub const ITALIC: &str = "\x1b[3m";
pub const UNDERLINE: &str = "\x1b[4m";
pub const BLINK: &str = "\x1b[5m";
pub const REVERSE: &str = "\x1b[7m";
pub const STRIKETHROUGH: &str = "\x1b[9m";

/// Lookup table: short code (uppercase) -> ANSI escape sequence.
pub fn colors() -> &'static HashMap<&'static str, &'static str> {
    static MAP: OnceLock<HashMap<&'static str, &'static str>> = OnceLock::new();
    MAP.get_or_init(|| {
        let mut m = HashMap::with_capacity(128);

        // Styles
        m.insert("BOLD", BOLD);
        m.insert("BD", BOLD);
        m.insert("DIM", DIM);
        m.insert("FAINT", DIM);
        m.insert("DM", DIM);
        m.insert("ITALIC", ITALIC);
        m.insert("IT", ITALIC);
        m.insert("UNDERLINE", UNDERLINE);
        m.insert("UE", UNDERLINE);
        m.insert("BLINK", BLINK);
        m.insert("BL", BLINK);
        m.insert("REVERSE", REVERSE);
        m.insert("RV", REVERSE);
        m.insert("STRIKETHROUGH", STRIKETHROUGH);
        m.insert("ST", STRIKETHROUGH);
        m.insert("RESET", RESET);
        m.insert("RT", RESET);

        // Foreground
        m.insert("BLACK", BLACK);
        m.insert("BK", BLACK);
        m.insert("RED", RED);
        m.insert("RD", RED);
        m.insert("GREEN", GREEN);
        m.insert("GN", GREEN);
        m.insert("YELLOW", YELLOW);
        m.insert("YW", YELLOW);
        m.insert("BLUE", BLUE);
        m.insert("BE", BLUE);
        m.insert("MAGENTA", MAGENTA);
        m.insert("MA", MAGENTA);
        m.insert("CYAN", CYAN);
        m.insert("CN", CYAN);
        m.insert("WHITE", WHITE);
        m.insert("WE", WHITE);

        // Bright foreground
        m.insert("BBLACK", BRIGHT_BLACK);
        m.insert("BBK", BRIGHT_BLACK);
        m.insert("BRED", BRIGHT_RED);
        m.insert("BRD", BRIGHT_RED);
        m.insert("BGREEN", BRIGHT_GREEN);
        m.insert("BGN", BRIGHT_GREEN);
        m.insert("BYELLOW", BRIGHT_YELLOW);
        m.insert("BYW", BRIGHT_YELLOW);
        m.insert("BBLUE", BRIGHT_BLUE);
        m.insert("BBE", BRIGHT_BLUE);
        m.insert("BMAGENTA", BRIGHT_MAGENTA);
        m.insert("BMA", BRIGHT_MAGENTA);
        m.insert("BCYAN", BRIGHT_CYAN);
        m.insert("BCN", BRIGHT_CYAN);
        m.insert("BWHITE", BRIGHT_WHITE);
        m.insert("BWE", BRIGHT_WHITE);

        // Extended foreground (xterm-256)
        m.insert("ORANGE", ORANGE);
        m.insert("OE", ORANGE);
        m.insert("PINK", PINK);
        m.insert("PK", PINK);
        m.insert("PURPLE", PURPLE);
        m.insert("PE", PURPLE);
        m.insert("VIOLET", VIOLET);
        m.insert("VT", VIOLET);
        m.insert("BROWN", BROWN);
        m.insert("BN", BROWN);
        m.insert("GRAY", GRAY);
        m.insert("GREY", GRAY);
        m.insert("GY", GRAY);
        m.insert("LIME", LIME);
        m.insert("LE", LIME);
        m.insert("TEAL", TEAL);
        m.insert("TL", TEAL);
        m.insert("NAVY", NAVY);
        m.insert("NY", NAVY);
        m.insert("GOLD", GOLD);
        m.insert("GD", GOLD);
        m.insert("SILVER", SILVER);
        m.insert("SR", SILVER);
        m.insert("MAROON", MAROON);
        m.insert("MN", MAROON);

        // Background standard
        m.insert("BACKBLACK", BLACK_BG);
        m.insert("BKBK", BLACK_BG);
        m.insert("BACKRED", RED_BG);
        m.insert("BKRD", RED_BG);
        m.insert("BACKGREEN", GREEN_BG);
        m.insert("BKGN", GREEN_BG);
        m.insert("BACKYELLOW", YELLOW_BG);
        m.insert("BKYW", YELLOW_BG);
        m.insert("BACKBLUE", BLUE_BG);
        m.insert("BKBE", BLUE_BG);
        m.insert("BACKMAGENTA", MAGENTA_BG);
        m.insert("BKMA", MAGENTA_BG);
        m.insert("BACKCYAN", CYAN_BG);
        m.insert("BKCN", CYAN_BG);
        m.insert("BACKWHITE", WHITE_BG);
        m.insert("BKWE", WHITE_BG);

        // Background bright
        m.insert("BACKBBLACK", BRIGHT_BLACK_BG);
        m.insert("BKBBK", BRIGHT_BLACK_BG);
        m.insert("BACKBRED", BRIGHT_RED_BG);
        m.insert("BKBRD", BRIGHT_RED_BG);
        m.insert("BACKBGREEN", BRIGHT_GREEN_BG);
        m.insert("BKBGN", BRIGHT_GREEN_BG);
        m.insert("BACKBYELLOW", BRIGHT_YELLOW_BG);
        m.insert("BKBYW", BRIGHT_YELLOW_BG);
        m.insert("BACKBBLUE", BRIGHT_BLUE_BG);
        m.insert("BKBBE", BRIGHT_BLUE_BG);
        m.insert("BACKBMAGENTA", BRIGHT_MAGENTA_BG);
        m.insert("BKBMA", BRIGHT_MAGENTA_BG);
        m.insert("BACKBCYAN", BRIGHT_CYAN_BG);
        m.insert("BKBCN", BRIGHT_CYAN_BG);
        m.insert("BACKBWHITE", BRIGHT_WHITE_BG);
        m.insert("BKBWE", BRIGHT_WHITE_BG);

        // Background extended (xterm-256)
        m.insert("BACKORANGE", ORANGE_BG);
        m.insert("BKOE", ORANGE_BG);
        m.insert("BACKPINK", PINK_BG);
        m.insert("BKPK", PINK_BG);
        m.insert("BACKPURPLE", PURPLE_BG);
        m.insert("BKPE", PURPLE_BG);
        m.insert("BACKVIOLET", VIOLET_BG);
        m.insert("BKVT", VIOLET_BG);
        m.insert("BACKBROWN", BROWN_BG);
        m.insert("BKBN", BROWN_BG);
        m.insert("BACKGRAY", GRAY_BG);
        m.insert("BACKGREY", GRAY_BG);
        m.insert("BKGY", GRAY_BG);
        m.insert("BACKLIME", LIME_BG);
        m.insert("BKLE", LIME_BG);
        m.insert("BACKTEAL", TEAL_BG);
        m.insert("BKTL", TEAL_BG);
        m.insert("BACKNAVY", NAVY_BG);
        m.insert("BKNY", NAVY_BG);
        m.insert("BACKGOLD", GOLD_BG);
        m.insert("BKGD", GOLD_BG);
        m.insert("BACKSILVER", SILVER_BG);
        m.insert("BKSR", SILVER_BG);
        m.insert("BACKMAROON", MAROON_BG);
        m.insert("BKMN", MAROON_BG);

        m
    })
}