use once_cell::sync::Lazy;
use std::path::PathBuf;
use std::env;
use std::fs;

pub const API_URL: &str = "http://localhost:8080";

pub static _EXECUTABLE_DIR: Lazy<PathBuf> = Lazy::new(|| {
    env::current_exe()
        .expect("failed to get current executable path")
        .parent()
        .expect("executable must have a parent directory")
        .to_path_buf()
});

pub static STAGE_PATH: Lazy<PathBuf> = Lazy::new(|| {
    let mut path = dirs::data_local_dir().expect("failed to get local data directory");
    path.push(".belchi/stage");

    if !path.exists() {
        println!("stage folder not found. Trying create");
        fs::create_dir_all(&path).expect("failed to create stage path");
        println!("stage folder has been created");
    } 

    path
});

pub static PROGRAM_PATH: Lazy<PathBuf> = Lazy::new(|| {
    let mut path = dirs::data_local_dir().expect("failed to get local data directory");
    path.push(".belchi/");

    if !path.exists() {
        println!("program folder not found. Trying create");
        fs::create_dir_all(&path).expect("failed to create program path");
        println!("program folder has been created");
    } 

    path
});

