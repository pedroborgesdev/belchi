use std::fs;
use std::path::Path;

use crate::config::STAGE_PATH;

pub fn run(file: &str) {
    let stage_path = STAGE_PATH.display().to_string();
    let stage_file_path = Path::new(&stage_path).join(format!("{}.belchi", file));

    println!("Deleting package '{}' from stage...", file);

    if !stage_file_path.exists() {
        println!("Failed to locate file in stage path.");
        println!("Reason: Package not found in the configured stage directory\n");
        return;
    }

    match fs::remove_file(&stage_file_path) {
        Ok(_) => {
            println!("\nPackage successfully deleted from stage.");
        }
        Err(_) => {
            println!("Failed to delete package from stage.");
            println!("Reason: Unable to remove specified file\n");
        }
    }
}
