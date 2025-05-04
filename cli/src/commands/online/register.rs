use std::io::{self, Write};
use reqwest;
use serde_json::json;

use crate::config::API_URL;

pub fn run() {
    println!("Registering account at {}...", API_URL);
    println!("Provide account information to create a new account securely.\n");

    let username = match prompt_input("Input account username: ") {
        Ok(u) => u,
        Err(_) => {
            eprintln!("\nFailed to read username.");
            return;
        }
    };

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
    let request_body = json!({
        "username": username,
        "email": email,
        "password": password
    });

    println!("\nSending registration data to the server...");

    match client.post(format!("{}/auth/register", API_URL))
        .json(&request_body)
        .send()
    {
        Ok(resp) => handle_response(resp),
        Err(_) => {
            eprintln!("\nFailed to send registration request.");
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
        println!("\nAccount created successfully!");
        if let Some(username) = body.get("data")
            .and_then(|data| data.get("username"))
            .and_then(|u| u.as_str()) 
        {
            println!("Welcome, {}!", username);
            println!("To login use: 'belchi login'\n");
        }
    } else {
        println!("\nRegistration failed!");
        if let Some(error) = body.get("data")
            .and_then(|data| data.get("error"))
            .and_then(|e| e.as_str())
        {
            println!("Reason: {}", error);
        }
    }
}

fn prompt_input(prompt: &str) -> io::Result<String> {
    print!("{}", prompt);
    io::stdout().flush()?;

    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    Ok(input.trim().to_string())
}
