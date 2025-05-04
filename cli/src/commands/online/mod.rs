use crate::cli::Command;

pub mod register;
pub mod login;

pub mod up;
pub mod get;

pub fn handle(command: Command) {
    match command {
        Command::Register => register::run(),
        Command::Login => login::run(),

        Command::Up(args) => up::run(&args.package),
        Command::Get(args) => get::run(&args.package),
        _ => unreachable!("Only online commands!"),
    }
}