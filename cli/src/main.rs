mod config;
mod cli;
mod commands;

use crate::config::STAGE_PATH;
use clap::Parser;

fn main() {
    let _ = &*STAGE_PATH;
    
    let args = cli::Cli::parse();
    commands::handle(args.command);
}
