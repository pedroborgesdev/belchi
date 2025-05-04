use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(
    bin_name = "belchi", 
    author = "Pedro Borges", 
    version = "1.0", 
    about = "Belchi - Structure Packager & Loader\nMaded by @pedroborgesdev\n\nContact me:\n  Linkedin:   Pedro Borges\n  Github:     pedroborgesdev", 
    long_about = None
)]
#[command(color = clap::ColorChoice::Never, after_help = "Enjoy Belchi!")]

pub struct Cli {
    #[command(subcommand)]
    pub command: Command,
}

#[derive(Subcommand)]
pub enum Command {
    #[command(about = "Register a new user")]
    Register,

    #[command(about = "Login as an existing user\n")]
    Login,

    #[command(about = "Record a new package structure")]
    Rec(RecArgs),

    #[command(about = "Save a file in stage with a specific version")]
    Stage(StageArgs),

    #[command(about = "Load a package by name")]
    Load(LoadArgs),

    #[command(about = "Move a package to a current location")]
    Move(MoveArgs),

    #[command(about = "Delete specified package at stage state")]
    Del(DelArgs),

    #[command(about = "Show all packages at stage state\n")]
    Show,

    #[command(about = "Upload a package")]
    Up(UpArgs),

    #[command(about = "Get a package from the server\n")]
    Get(GetArgs),
}

#[derive(clap::Args)]
pub struct RecArgs {
    #[arg(help = "Name of the new structure to record")]
    pub name: String,
}

#[derive(clap::Args)]
pub struct LoadArgs {
    #[arg(help = "Name of the package to load")]
    pub package: String,
}

#[derive(clap::Args)]
pub struct MoveArgs {
    #[arg(help = "Name of the package to move")]
    pub package: String,
}

#[derive(clap::Args)]
pub struct UpArgs {
    #[arg(help = "Name of the package to upload")]
    pub package: String,
}

#[derive(clap::Args)]
pub struct StageArgs {
    #[arg(help = "Name of the file to save")]
    pub file: String,

    #[arg(help = "Version tag to associate with the file")]
    pub version: String,
}

#[derive(clap::Args)]
pub struct GetArgs {
    #[arg(help = "Name of the package to download")]
    pub package: String,
}

#[derive(clap::Args)]
pub struct DelArgs {
    #[arg(help = "Name of the package to remove")]
    pub package: String,
}
