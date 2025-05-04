use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

use crate::config::STAGE_PATH;

pub fn run(name: &str) {
    println!("Loading package '{}' from stage...", name);

    let script_path = STAGE_PATH.join(format!("{}.belchi", name));

    let content = match fs::read_to_string(&script_path) {
        Ok(c) => c,
        Err(_) => {
            eprintln!("Failed to read specified file.");
            println!("Reason: package not found on stage.\n");
            return;
        }
    };

    let mut folders = HashMap::new();
    let mut files = Vec::new();
    let mut contents = HashMap::new();
    let mut current_section = None;

    for line in content.lines().map(str::trim) {
        match line {
            "!FOLDERS [" => current_section = Some("FOLDERS"),
            "!FILES [" => current_section = Some("FILES"),
            "!CONTENT [" => current_section = Some("CONTENT"),
            "]" => current_section = None,
            _ => {
                if let Some(section) = current_section {
                    match section {
                        "FOLDERS" => {
                            if let Some((idx, path)) = line.split_once('=') {
                                if let Ok(index) = idx.trim().parse::<usize>() {
                                    folders.insert(index, path.trim().to_string());
                                }
                            }
                        }
                        "FILES" => {
                            if let Some((file_idx, rest)) = line.split_once('>') {
                                if let (Ok(file_index), Some((folder_idx, name))) =
                                    (file_idx.trim().parse(), rest.split_once('='))
                                {
                                    if let Ok(folder_index) = folder_idx.trim().parse::<usize>() {
                                        files.push((file_index, folder_index, name.trim().to_string()));
                                    }
                                }
                            }
                        }
                        "CONTENT" => {
                            if let Some((file_idx, content)) = line.split_once("> --") {
                                if let Ok(file_index) = file_idx.trim().parse::<usize>() {
                                    let cleaned = content.trim_end_matches("--").replace("\\n", "\n");
                                    contents.insert(file_index, cleaned);
                                }
                            }
                        }
                        _ => {}
                    }
                }
            }
        }
    }

    println!("\nCreating required folders...");

    for path in folders.values().filter(|p| !p.is_empty()) {
        if let Err(_) = fs::create_dir_all(path) {
            eprintln!("\nFailed to create one or more folders.");
            println!("Reason: unable to create directory structure.\n");
            return;
        }
    }

    println!("Creating files and applying contents...\n");

    let mut created = 0;

    for (file_idx, folder_idx, filename) in files {
        let folder_path = folders.get(&folder_idx).cloned().unwrap_or_default();
        let full_path = if folder_path.is_empty() {
            PathBuf::from(&filename)
        } else {
            Path::new(&folder_path).join(&filename)
        };

        if let Some(parent) = full_path.parent() {
            if let Err(_) = fs::create_dir_all(parent) {
                eprintln!("Failed to prepare folder structure for file.");
                println!("Reason: unable to create parent directory.\n");
                return;
            }
        }

        let file_content = contents.get(&file_idx).cloned().unwrap_or_default();
        match fs::write(&full_path, file_content) {
            Ok(_) => created += 1,
            Err(_) => {
                eprintln!("Failed to write file.");
                println!("Reason: could not write file content to disk.\n");
                return;
            }
        }
    }

    if created == 0 {
        println!("Package could not be loaded.\n");
        println!("Reason: invalid or empty package.\n");
    } else {
        println!("Package loaded successfully!");
    }
}
