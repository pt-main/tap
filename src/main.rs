use tap::{engine::structs::Engine, utils::add_help_handler};

fn main() {
    let mut ap = tap::argsparsing::ArgsP::new();
    let mut eng = Engine::new(&mut ap, "".to_string());
    add_help_handler(vec!["help", "h"], &mut eng);
    if let Err(e) = eng.main() {
        println!("{}", e);
    }
}