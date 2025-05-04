use std::fs::{self, OpenOptions};
use std::io::{self, Write, BufReader, BufRead};
use std::path::Path;

use crate::config::STAGE_PATH;

pub fn run(file: &str, version: &str) {
    let stage_path = STAGE_PATH.display().to_string();

    println!("Staging package '{}'", file);
    println!("Updating version to '{}'\n", version);

    if let Err(e) = moving_file(file, &stage_path, version) {
        println!("\nFailed to stage the package.");
        println!("Reason: {}\n", e);
    }
}

fn moving_file(file: &str, path: &str, version: &str) -> io::Result<()> {
    let file_path = Path::new(file).with_extension("belchi");
    let target_path = Path::new(path).join(file_path.file_name().unwrap());

    if !file_path.exists() {
        println!("\nFailed to locate specified file.");
        println!("Reason: '{}.belchi' not found\n", file);
        return Ok(());
    }

    if !Path::new(path).exists() {
        println!("Stage folder not found. Creating...\n");
        fs::create_dir_all(path)?;
        println!("Stage folder created at '{}'\n", path);
    }

    if target_path.exists() {
        print!("File already exists. Do you want to overwrite it? (y/N): ");
        io::stdout().flush()?;

        let mut input = String::new();
        io::stdin().read_line(&mut input)?;
        let input = input.trim().to_lowercase();

        if input != "y" {
            println!("\nOperation aborted. File was not overwritten.\n");
            return Ok(());
        }
    }

    let file = OpenOptions::new().read(true).open(&file_path)?;
    let reader = BufReader::new(file);
    let mut lines: Vec<String> = reader.lines().filter_map(Result::ok).collect();

    let mut version_found = false;
    for line in lines.iter_mut() {
        if line.starts_with("!VERSION=") {
            *line = format!("!VERSION={}", version);
            version_found = true;
            break;
        }
    }

    if !version_found {
        lines.push(format!("\n!VERSION={}", version));
    }

    println!("Copying file to stage folder...");
    fs::copy(&file_path, &target_path)?;

    let mut writer = OpenOptions::new().write(true).truncate(true).open(&target_path)?;
    for line in lines {
        writeln!(writer, "{}", line)?;
    }

    println!("\nPackage successfully added to stage.");

    Ok(())
}
