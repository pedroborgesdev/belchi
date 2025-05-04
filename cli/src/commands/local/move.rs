use std::fs;
use std::io::{self, Write};
use std::path::Path;

use crate::config::STAGE_PATH;

pub fn run(file: &str) {
    let stage_path = STAGE_PATH.display().to_string();
    let stage_file_path = Path::new(&stage_path).join(format!("{}.belchi", file));
    let target_file_name = format!("{}.belchi", file);
    let target_file_path = Path::new(&target_file_name);

    println!("Retrieving package '{}' from stage...", file);

    if !stage_file_path.exists() {
        println!("\nFailed to locate file in stage path.");
        println!("Reason: Package not found in the configured stage directory\n");
        return;
    }

    if target_file_path.exists() {
        print!("\nFile existing in current directory. Overwrite it? (y/N): ");
        io::stdout().flush().unwrap();

        let mut input = String::new();
        if let Err(_) = io::stdin().read_line(&mut input) {
            println!("\nFailed to read user input.");
            println!("Reason: Unable to access standard input\n");
            return;
        }

        if input.trim().to_lowercase() != "y" {
            println!("\nOperation aborted by user.");
            println!("No file was overwritten\n");
            return;
        }
    }

    println!("Copying file from stage...\n");

    match fs::copy(&stage_file_path, &target_file_path) {
        Ok(_) => {
            println!("Package file successfully retrieved from stage.");
        }
        Err(_) => {
            println!("Failed to retrieve the package from stage.");
            println!("Reason: Unable to copy specified file\n");
        }
    }
}
