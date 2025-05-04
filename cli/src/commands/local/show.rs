use std::fs;
use std::path::Path;

use crate::config::STAGE_PATH;

pub fn run() {
    let stage_path = STAGE_PATH.display().to_string();
    let path = Path::new(&stage_path);

    println!("Listing staged packages...\n");

    if !path.exists() {
        println!("Failed to access stage directory.");
        println!("Reason: Configured stage path does not exist\n");
        return;
    }

    let entries = match fs::read_dir(path) {
        Ok(e) => e,
        Err(_) => {
            println!("Failed to read stage directory.");
            println!("Reason: Unable to access directory contents\n");
            return;
        }
    };

    let mut count = 0;

    for entry in entries {
        match entry {
            Ok(entry) => {
                let file_path = entry.path();

                if file_path.extension().and_then(|ext| ext.to_str()) == Some("belchi") {
                    if let Some(file_name) = file_path.file_stem() {
                        println!("- {}", file_name.to_string_lossy());
                        count += 1;
                    }
                }
            }
            Err(_) => {
                println!("Failed to process a file entry.");
                println!("Reason: Entry could not be read\n");
            }
        }
    }

    if count == 0 {
        println!("No packages found in stage.");
        println!("Ensure that package files are correctly staged before listing\n");
    } else {
        println!("\nTotal packages found: {}", count);
    }
}
