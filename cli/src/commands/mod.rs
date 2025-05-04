pub mod online;
pub mod local;

use crate::cli::Command;

pub fn handle(command: Command) {
    match command {
        Command::Login | Command::Register | Command::Up(_) | Command::Get(_) => online::handle(command),
        Command::Rec(_) | Command::Load(_) | Command::Stage(_) | Command::Move(_) | Command::Show | Command::Del(_)  => local::handle(command),    
    }
}