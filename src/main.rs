use std::process::Command;
use std::io::{self, Write};

fn adb(cmd: &str) -> String {
    let output = Command::new("adb")
        .args(&["shell", cmd])
        .output()
        .expect("Failed to run adb");
    String::from_utf8_lossy(&output.stdout).to_string()
}

fn list_packages(filter: Option<&str>, user_only: bool) -> Vec<String> {
    let flag = if user_only { "-3" } else { "" };
    let cmd = format!("pm list packages {}", flag);
    adb(&cmd)
        .lines()
        .filter_map(|line| {
            let pkg = line.strip_prefix("package:")?.to_string();
            if let Some(f) = filter {
                if pkg.contains(f) { Some(pkg) } else { None }
            } else {
                Some(pkg)
            }
        })
        .collect()
}

fn get_app_size(pkg: &str) -> String {
    let out = adb(&format!("du -sh /data/app/{}", pkg));
    out.lines().next().unwrap_or("?").to_string()
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    
    if args.len() < 2 {
        println!("Usage:");
        println!("  android-pm list                 -- list all packages");
        println!("  android-pm list -u              -- user apps only");
        println!("  android-pm search <keyword>     -- search packages");
        println!("  android-pm size <package>       -- show package size");
        return;
    }

    match args[1].as_str() {
        "list" => {
            let user_only = args.contains(&"-u".to_string());
            let pkgs = list_packages(None, user_only);
            println!("Found {} packages\n", pkgs.len());
            for (i, pkg) in pkgs.iter().enumerate() {
                if i % 2 == 0 {
                    print!("{:<40}", pkg);
                } else {
                    println!("{}", pkg);
                }
            }
            if pkgs.len() % 2 != 0 {
                println!();
            }
        },
        "search" => {
            if args.len() < 3 {
                eprintln!("Usage: android-pm search <keyword>");
                return;
            }
            let pkgs = list_packages(Some(&args[2]), false);
            println!("Found {} matches:\n", pkgs.len());
            for pkg in pkgs {
                println!("  {}", pkg);
            }
        },
        "size" => {
            if args.len() < 3 {
                eprintln!("Usage: android-pm size <package>");
                return;
            }
            let size = get_app_size(&args[2]);
            println!("Size: {}", size);
        },
        _ => println!("Unknown command: {}", args[1]),
    }
}
