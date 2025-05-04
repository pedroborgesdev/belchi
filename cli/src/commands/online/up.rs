use reqwest::blocking::multipart;
use std::fs;

use crate::config::{API_URL, STAGE_PATH, PROGRAM_PATH};

pub fn run(package: &str) {
    println!("Uploading package to {}", API_URL);

    let script_path = STAGE_PATH.join(format!("{}.belchi", package));
    let content = match fs::read_to_string(&script_path) {
        Ok(c) => {
            c
        }
        Err(_) => {
            eprintln!("\nFailed to read specified file");
            println!("Reason: Package not found on stage.");
            return;
        }
    };

    let mut version = "1".to_string();
    let mut name = "undefined".to_string();

    for line in content.lines().map(str::trim) {
        if let Some(ver) = line.strip_prefix("!VERSION=") {
            version = ver.trim().to_string();
        } else if let Some(n) = line.strip_prefix("!NAME=") {
            name = n.trim().to_string();
        }
    }

    let form = match multipart::Form::new()
        .text("name", name.clone())
        .text("version", version.clone())
        .file("file", &script_path)
    {
        Ok(form) => {
            form
        }
        Err(_) => {
            eprintln!("\nCannot create form to upload");
            return;
        }
    };

    let client = reqwest::blocking::Client::new();


    let token_path = PROGRAM_PATH.join("auth_data.txt");
    let content = match fs::read_to_string(&token_path) {
        Ok(c) => {
            c
        }
        Err(_) => {
            eprintln!("\nFailed to read specified file");
            println!("Reason: Unable to read authentication data.");
            return;
        }
    };

    let mut token = "undefined".to_string();
    for line in content.lines().map(str::trim) {
        if let Some(tk) = line.strip_prefix("token=") {
            token = tk.trim().to_string();
        }
    }

    if token == "undefined" {
        eprintln!("\nFailed to get authorization token");
        println!("Reason: Token not found. Try logging in again.\n");
        return;
    }

    println!("Sending package to server...");

    match client
        .post(format!("{}/packages/", API_URL))
        .multipart(form)
        .bearer_auth(token)
        .send()
    {
        Ok(resp) => handle_response(resp),
        Err(_) => {
            eprintln!("\nError occurred during request");
            println!("Reason: Unable to connect or send request to server.\n");
        }
    }
}

fn handle_response(resp: reqwest::blocking::Response) {
    let status = resp.status();

    let body: serde_json::Value = match resp.json() {
        Ok(b) => b,
        Err(_) => {
            eprintln!("\nFailed to parse server response");
            return;
        }
    };

    if status.is_success() {
        println!("\nPackage uploaded successfully!");
        if let Some(package) = body.get("data")
            .and_then(|data| data.get("package"))
            .and_then(|e| e.as_str())
        {
            println!("To retrieve your package, use 'belchi get {}'\n", package);
        }
    } else {
        println!("\nPackage upload failed!");
        if let Some(error) = body.get("data")
            .and_then(|data| data.get("error"))
            .and_then(|e| e.as_str())
        {
            println!("Reason: {}\n", error);
        }
    }
}
