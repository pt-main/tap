use std::{collections::HashMap, sync::OnceLock};

use crate::tooling::console::styles::{Color, Str, Text};

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


macro_rules! add_color {
    ($m:expr, $name:literal, $code:expr; $($alias:literal),+ $(,)?) => {
        $(
            $m.insert($alias, Color { name: $name, code: $code });
        )+
    };
}

pub fn colors() -> &'static HashMap<&'static str, Color> {
    static MAP: OnceLock<HashMap<&'static str, Color>> = OnceLock::new();
    MAP.get_or_init(|| {
        let mut m = HashMap::with_capacity(128);

        // --- Styles ---
        add_color!(m, "bold",          BOLD;          "BOLD", "BD");
        add_color!(m, "dim",           DIM;           "DIM", "FAINT", "DM");
        add_color!(m, "italic",        ITALIC;        "ITALIC", "IT");
        add_color!(m, "underline",     UNDERLINE;     "UNDERLINE", "UE");
        add_color!(m, "blink",         BLINK;         "BLINK", "BL");
        add_color!(m, "reverse",       REVERSE;       "REVERSE", "RV");
        add_color!(m, "strikethrough", STRIKETHROUGH; "STRIKETHROUGH", "ST");
        add_color!(m, "reset",         RESET;         "RESET", "RT");

        // --- Foreground (standard) ---
        add_color!(m, "black",   BLACK;   "BLACK", "BK");
        add_color!(m, "red",     RED;     "RED", "RD");
        add_color!(m, "green",   GREEN;   "GREEN", "GN");
        add_color!(m, "yellow",  YELLOW;  "YELLOW", "YW");
        add_color!(m, "blue",    BLUE;    "BLUE", "BE");
        add_color!(m, "magenta", MAGENTA; "MAGENTA", "MA");
        add_color!(m, "cyan",    CYAN;    "CYAN", "CN");
        add_color!(m, "white",   WHITE;   "WHITE", "WE");

        // --- Foreground (bright) ---
        add_color!(m, "bright_black",   BRIGHT_BLACK;   "BBLACK", "BBK");
        add_color!(m, "bright_red",     BRIGHT_RED;     "BRED", "BRD");
        add_color!(m, "bright_green",   BRIGHT_GREEN;   "BGREEN", "BGN");
        add_color!(m, "bright_yellow",  BRIGHT_YELLOW;  "BYELLOW", "BYW");
        add_color!(m, "bright_blue",    BRIGHT_BLUE;    "BBLUE", "BBE");
        add_color!(m, "bright_magenta", BRIGHT_MAGENTA; "BMAGENTA", "BMA");
        add_color!(m, "bright_cyan",    BRIGHT_CYAN;    "BCYAN", "BCN");
        add_color!(m, "bright_white",   BRIGHT_WHITE;   "BWHITE", "BWE");

        // --- Foreground (extended xterm-256) ---
        add_color!(m, "orange", ORANGE; "ORANGE", "OE");
        add_color!(m, "pink",   PINK;   "PINK", "PK");
        add_color!(m, "purple", PURPLE; "PURPLE", "PE");
        add_color!(m, "violet", VIOLET; "VIOLET", "VT");
        add_color!(m, "brown",  BROWN;  "BROWN", "BN");
        add_color!(m, "gray",   GRAY;   "GRAY", "GREY", "GY");
        add_color!(m, "lime",   LIME;   "LIME", "LE");
        add_color!(m, "teal",   TEAL;   "TEAL", "TL");
        add_color!(m, "navy",   NAVY;   "NAVY", "NY");
        add_color!(m, "gold",   GOLD;   "GOLD", "GD");
        add_color!(m, "silver", SILVER; "SILVER", "SR");
        add_color!(m, "maroon", MAROON; "MAROON", "MN");

        // --- Background (standard) ---
        add_color!(m, "bg_black",   BLACK_BG;   "BACKBLACK", "BKBK");
        add_color!(m, "bg_red",     RED_BG;     "BACKRED", "BKRD");
        add_color!(m, "bg_green",   GREEN_BG;   "BACKGREEN", "BKGN");
        add_color!(m, "bg_yellow",  YELLOW_BG;  "BACKYELLOW", "BKYW");
        add_color!(m, "bg_blue",    BLUE_BG;    "BACKBLUE", "BKBE");
        add_color!(m, "bg_magenta", MAGENTA_BG; "BACKMAGENTA", "BKMA");
        add_color!(m, "bg_cyan",    CYAN_BG;    "BACKCYAN", "BKCN");
        add_color!(m, "bg_white",   WHITE_BG;   "BACKWHITE", "BKWE");

        // --- Background (bright) ---
        add_color!(m, "bg_bright_black",   BRIGHT_BLACK_BG;   "BACKBBLACK", "BKBBK");
        add_color!(m, "bg_bright_red",     BRIGHT_RED_BG;     "BACKBRED", "BKBRD");
        add_color!(m, "bg_bright_green",   BRIGHT_GREEN_BG;   "BACKBGREEN", "BKBGN");
        add_color!(m, "bg_bright_yellow",  BRIGHT_YELLOW_BG;  "BACKBYELLOW", "BKBYW");
        add_color!(m, "bg_bright_blue",    BRIGHT_BLUE_BG;    "BACKBBLUE", "BKBBE");
        add_color!(m, "bg_bright_magenta", BRIGHT_MAGENTA_BG; "BACKBMAGENTA", "BKBMA");
        add_color!(m, "bg_bright_cyan",    BRIGHT_CYAN_BG;    "BACKBCYAN", "BKBCN");
        add_color!(m, "bg_bright_white",   BRIGHT_WHITE_BG;   "BACKBWHITE", "BKBWE");

        // --- Background (extended) ---
        add_color!(m, "bg_orange", ORANGE_BG; "BACKORANGE", "BKOE");
        add_color!(m, "bg_pink",   PINK_BG;   "BACKPINK", "BKPK");
        add_color!(m, "bg_purple", PURPLE_BG; "BACKPURPLE", "BKPE");
        add_color!(m, "bg_violet", VIOLET_BG; "BACKVIOLET", "BKVT");
        add_color!(m, "bg_brown",  BROWN_BG;  "BACKBROWN", "BKBN");
        add_color!(m, "bg_gray",   GRAY_BG;   "BACKGRAY", "BACKGREY", "BKGY");
        add_color!(m, "bg_lime",   LIME_BG;   "BACKLIME", "BKLE");
        add_color!(m, "bg_teal",   TEAL_BG;   "BACKTEAL", "BKTL");
        add_color!(m, "bg_navy",   NAVY_BG;   "BACKNAVY", "BKNY");
        add_color!(m, "bg_gold",   GOLD_BG;   "BACKGOLD", "BKGD");
        add_color!(m, "bg_silver", SILVER_BG; "BACKSILVER", "BKSR");
        add_color!(m, "bg_maroon", MAROON_BG; "BACKMAROON", "BKMN");

        m
    })
}

pub fn colorize(inp: &str, text: &mut Text) {
    let mut skip: i32 = -1;
    for (idx, ch) in inp.char_indices() {
        if skip > (idx as i32) {
            continue;
        }
        let curr = &inp[idx..];
        let mut push = true;
        if ch == '[' {
            if let Some(new) = curr.strip_prefix("[?") {
                let mut best: Option<(&str, Color)> = None;
                for (code, col) in colors() {
                    if new.starts_with(code) {
                        if best.map_or(true, |(b, _)| code.len() > b.len()) {
                            best = Some((code, *col));
                        }
                    }
                }
                if let Some((code, col)) = best {
                    let mut s = Str::new();
                    s.colors.add(col);
                    text.text.push(s);
                    push = false;
                    skip = idx as i32 + code.len() as i32 + 3;
                }
            }
        }
        if push {
            let l = text.text.len();
            if l > 0 {
                text.text[l-1].text.push(ch);
            } else {
                let mut str = Str::new();
                str.text.push(ch);
                text.text.push(str);
            }
        }
    }
}