use std::collections::HashMap;
use std::fs;
use std::io;
use std::path::Path;

#[derive(Debug)]
struct ProjectStructure {
    name: String,
    folders: Vec<(usize, String)>,
    files: Vec<(usize, usize, String)>,
    contents: Vec<(usize, String)>,
    ignore_dirs: Vec<String>,
    ignore_exts: Vec<String>,
}

impl ProjectStructure {
    fn new(name: &str) -> Self {
        Self {
            name: name.into(),
            folders: Vec::new(),
            files: Vec::new(),
            contents: Vec::new(),
            ignore_dirs: Vec::new(),
            ignore_exts: Vec::new(),
        }
    }

    fn add_ignored_dirs(&mut self, dirs: Vec<&str>) {
        self.ignore_dirs = dirs.into_iter().map(String::from).collect();
    }

    fn add_ignored_exts(&mut self, exts: Vec<&str>) {
        self.ignore_exts = exts.into_iter().map(|e| e.to_lowercase()).collect();
    }

    fn should_ignore(&self, path: &Path) -> bool {
        let path_str = path.to_string_lossy();
        self.ignore_dirs.iter().any(|d| path_str.contains(d))
            || path
                .extension()
                .and_then(|ext| ext.to_str())
                .map_or(false, |ext| self.ignore_exts.contains(&ext.to_lowercase()))
    }

    fn scan_directory(&mut self, base_path: &Path) -> io::Result<()> {
        println!("Scanning project directory structure");

        let mut folder_map = HashMap::new();
        let mut file_counter = 0;

        self.index_folders(base_path, base_path, &mut folder_map)?;

        for entry in walkdir::WalkDir::new(base_path)
            .into_iter()
            .filter_map(Result::ok)
        {
            let path = entry.path();

            if self.should_ignore(path) || !entry.file_type().is_file() {
                continue;
            }

            let parent_path = path.parent().unwrap_or(base_path);
            let parent_key = parent_path
                .strip_prefix(base_path)
                .unwrap_or(parent_path)
                .to_string_lossy()
                .to_string();

            let folder_index = *folder_map.get(&parent_key).unwrap_or(&0);
            let file_name = path.file_name().unwrap().to_string_lossy().into_owned();

            self.files.push((file_counter, folder_index, file_name));

            let content = match fs::read_to_string(path) {
                Ok(c) => c,
                Err(_) => {
                    eprintln!("Failed to read specified file");
                    continue;
                }
            };

            self.contents.push((file_counter, content));
            file_counter += 1;
        }

        Ok(())
    }

    fn index_folders(
        &mut self,
        base: &Path,
        current: &Path,
        map: &mut HashMap<String, usize>,
    ) -> io::Result<()> {
        if map.is_empty() {
            self.folders.push((0, "./".to_string()));
            map.insert("".to_string(), 0);
        }

        if self.should_ignore(current) {
            return Ok(());
        }

        let relative = current.strip_prefix(base).unwrap_or(current);
        let key = relative.to_string_lossy();

        if !key.is_empty() && !map.contains_key(key.as_ref()) {
            let index = self.folders.len();
            self.folders.push((index, key.to_string()));
            map.insert(key.to_string(), index);
        }

        if current.is_dir() {
            for entry in fs::read_dir(current)? {
                let entry = entry?;
                if entry.file_type()?.is_dir() {
                    self.index_folders(base, &entry.path(), map)?;
                }
            }
        }

        Ok(())
    }

    fn serialize(&self) -> String {
        let mut output = format!("!NAME={}\n!VERSION=\n", self.name);

        output += "!FOLDERS [\n";
        for (i, p) in &self.folders {
            output += &format!("    {} = {}\n", i, p);
        }
        output += "]\n";

        output += "!FILES [\n";
        for (f, d, n) in &self.files {
            output += &format!("    {} > {} = {}\n", f, d, n);
        }
        output += "]\n";

        output += "!CONTENT [\n";
        for (f, c) in &self.contents {
            let escaped = c.replace('\n', "\\n");
            output += &format!("    {} > --{}--\n", f, escaped);
        }
        output += "]\n";

        output
    }
}

pub fn run(name: &str) {
    println!("Starting project structure generation for '{}'", name);

    if name.len() < 6 {
        eprintln!("\nStructure cannot be saved!");
        eprintln!("Reason: name must be at least 6 characters.");
        return;
    }

    let mut project = ProjectStructure::new(name);

    project.add_ignored_exts(vec![
        "exe", "dll", "class", "pyc", "o", "obj", "lock", "tmp", "log", "bin", "so", "a",
        "db", "sqlite3", "zip", "tar", "gz", "7z", "rar", "swp", "swo", "bak", "tmp", "jar",
    ]);
    project.add_ignored_dirs(vec![
        "target", ".git", "node_modules", "build", "dist", "out", ".idea", ".vscode",
        "__pycache__", ".cache", ".pytest_cache", ".mypy_cache", "venv", "env", ".env",
        ".DS_Store", "logs", "tmp", "coverage", ".gradle", ".settings",
    ]);

    if let Err(_) = project.scan_directory(Path::new(".")) {
        eprintln!("\nFailed to scan project directory");
        return;
    }

    let file_name = format!("{}.belchi", name);
    match fs::write(&file_name, project.serialize()) {
        Ok(_) => {
            println!("\nStructure saved to '{}'", file_name);
        }
        Err(_) => {
            eprintln!("\nFailed to save structure!");
            eprintln!("Reason: could not write output file");
        }
    }
}
