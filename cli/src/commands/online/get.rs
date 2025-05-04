use std::fs::File;
use std::io::{self, copy, Write};
use reqwest::header::CONTENT_DISPOSITION;

use crate::config::{API_URL, STAGE_PATH};

pub fn run(package: &str) {
    println!("Fetching package '{}' from {}...", package, API_URL);

    let client = reqwest::blocking::Client::new();
    let url = format!("{}/packages/{}", API_URL, package);

    match client.get(&url).send() {
        Ok(resp) => {
            handle_response(resp);
        }
        Err(_) => {
            eprintln!("\nFailed to send request to the server.");
            println!("Reason: Could not connect to the specified endpoint.\n");
        }
    }
}

fn handle_response(mut resp: reqwest::blocking::Response) {
    println!("Processing server response...");

    let status = resp.status();

    if !status.is_success() {
        eprintln!("\nFailed to download package.");
        println!("Reason: Package not exists.\n");
        return;
    }

    let mut filename = String::from("downloaded_file");

    if let Some(header_value) = resp.headers().get(CONTENT_DISPOSITION) {
        if let Ok(value_str) = header_value.to_str() {
            if let Some(name) = value_str.split(';')
                .find_map(|part| {
                    let part = part.trim();
                    if part.starts_with("filename=") {
                        Some(part.trim_start_matches("filename=").trim_matches('"'))
                    } else {
                        None
                    }
                }) {
                filename = name.to_string();
            }
        }
    }

    let path = STAGE_PATH.join(&filename);

    if path.exists() {
        print!("\nFile already exists. Overwrite it? (y/N): ");
        io::stdout().flush().expect("Failed to flush stdout");

        let mut input = String::new();
        io::stdin().read_line(&mut input).expect("Failed to read input");
        let input = input.trim().to_lowercase();

        if input != "y" {
            println!("\nOperation aborted. File was not overwritten.\n");
            return;
        }
    }

    println!("Saving package as in staging directory...");

    match File::create(&path) {
        Ok(mut file) => {
            println!("Writing contents to file...");

            if let Err(_) = copy(&mut resp, &mut file) {
                eprintln!("\nFailed to write file.");
                println!("Reason: Could not write data to disk.\n");
            } else {
                println!("\nPackage successfully saved.");
            }
        }
        Err(_) => {
            eprintln!("\nFailed to create file.");
            println!("Reason: Could not create the specified output file.\n");
        }
    }
}