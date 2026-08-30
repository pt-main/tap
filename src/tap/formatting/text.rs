use crate::formatting::set::colorize;

pub fn frame(text: &str, prefix: &str, postfix: &str, color: &str) -> String {
    let mut res= vec![format!("[?{}]╭─────── [?RT]{}", color, prefix)];
    for line in text.split("\n") {
        if line.starts_with("-") {
            res.push(format!("[?{}]⎬─[?RT]{}", color, line.replacen("-", "", 1)));
        } else {
            res.push(format!("[?{}]│    [?RT]{}", color, line));
        }
    }
    res.push(format!("[?{}]╰─────── [?RT]{}", color, postfix));
    colorize(res.join("\n").as_str())
}

pub fn add_space (text: &str, space: &str, n_spaces: usize) -> String {
    let space = &space.repeat(n_spaces);
    let mut res = String::from(space);
    res.push_str(&text.replace("\n", space));
    res
}