use std::process::Command;
use std::path::PathBuf;
use std::env::consts::{OS, ARCH};

fn get_binary_path() -> PathBuf {
    let os_name = match OS {
        "windows" => "win",
        "macos"   => "darwin",
        other     => other,
    };

    let arch_name = match ARCH {
        "x86_64"  => "amd64",
        "aarch64" => "arm64",
        other     => other,
    };

    let mut binary_name = format!("cmd-{}-{}-final", os_name, arch_name);
    if OS == "windows" {
        binary_name.push_str(".exe");
    }
    PathBuf::from("./resources").join(binary_name)
}

pub fn set(text: &str) -> String {
    let binary_path = get_binary_path();
    if !binary_path.exists() {
        panic!("Binary not found: {:?}", binary_path);
    }
    let output = Command::new(&binary_path)
        .arg(format!("\"{}\"", text.replace("\"", "\\\""))) 
        .output()
        .expect("failed to execute process");

    let res = String::from_utf8(output.stdout)
        .expect("output is not valid UTF-8");
    res[1..(res.len()-1)].to_string()
}