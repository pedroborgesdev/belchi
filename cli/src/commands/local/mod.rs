use crate::cli::Command;

pub mod rec;
pub mod load;
pub mod stage;
pub mod r#move;
pub mod show;
pub mod del;

pub fn handle(command: Command) {
    match command {
        Command::Rec(args) => rec::run(&args.name),
        Command::Load(args) => load::run(&args.package),
        Command::Stage(args) => stage::run(&args.file, &args.version),
        Command::Show => show::run(),
        Command::Del(args) => del::run(&args.package),
        Command::Move(args) => r#move::run(&args.package),
        _ => unreachable!("Only local commands!"),
    }
}