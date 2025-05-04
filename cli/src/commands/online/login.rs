use std::io::{self, Write};
use std::fs;
use reqwest;
use serde_json::json;

use crate::config::{API_URL, PROGRAM_PATH};

pub fn run() {
    println!("Logging in to {}...", API_URL);
    println!("Provide your credentials to receive an authentication token.\n");

    let email = match prompt_input("Input account email: ") {
        Ok(e) => e,
        Err(_) => {
            eprintln!("\nFailed to read email.");
            return;
        }
    };

    let password = match prompt_input("Input account password: ") {
        Ok(p) => p,
        Err(_) => {
            eprintln!("\nFailed to read password.");
            return;
        }
    };

    let client = reqwest::blocking::Client::new();
    let request_body = json!({ "email": email, "password": password });

    println!("\nSending credentials to the server...");

    match client.post(format!("{}/auth/login", API_URL))
        .json(&request_body)
        .send() 
    {
        Ok(resp) => handle_response(resp),
        Err(_) => {
            eprintln!("\nFailed to send login request.");
            println!("Reason: Unable to reach the server. Check your internet connection or try again later.");
        }
    }
}

fn handle_response(resp: reqwest::blocking::Response) {
    let status = resp.status();
    let body: serde_json::Value = match resp.json() {
        Ok(body) => body,
        Err(_) => {
            eprintln!("Failed to parse server response.");
            return;
        }
    };

    if status.is_success() {
        let token = body.get("data")
            .and_then(|data| data.get("token"))
            .and_then(|t| t.as_str());

        let username = body.get("data")
            .and_then(|data| data.get("username"))
            .and_then(|u| u.as_str());

        match (token, username) {
            (Some(token), Some(username)) => {
                println!("Saving authentication token locally...");

                if save_auth_data(token) {
                    println!("\nWelcome, {}!", username);
                    println!("Your token was saved successfully.");
                } else {
                    eprintln!("\nFailed to save token to file. Please try logging in again.");
                }
            }
            _ => eprintln!("Token or username not found in server response."),
        }
    } else {
        println!("\nLogin failed!");
        if let Some(error) = body.get("data")
            .and_then(|data| data.get("error"))
            .and_then(|e| e.as_str())
        {
            println!("Reason: {}", error);
        }
    }
}

fn save_auth_data(token: &str) -> bool {
    fs::write(
        format!("{}/auth_data.txt", PROGRAM_PATH.display()),
        format!("token={}", token)
    ).is_ok()
}

fn prompt_input(prompt: &str) -> io::Result<String> {
    print!("{}", prompt);
    io::stdout().flush()?;

    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    Ok(input.trim().to_string())
}
