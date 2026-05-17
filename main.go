package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strings"
    "sync"
    "time"
)

type Package struct {
    Name        string `json:"name"`
    Size        int64  `json:"size"`
    InstallTime int64  `json:"install_time"`
    IsSystem    bool   `json:"is_system"`
    Category    string `json:"category"`
}

func adb(args ...string) (string, error) {
    cmd := exec.Command("adb", args...)
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}

func listPackages(userOnly bool) ([]Package, error) {
    flag := ""
    if userOnly {
        flag = "-3"
    }
    
    out, err := adb("shell", fmt.Sprintf("pm list packages %s", flag))
    if err != nil {
        return nil, err
    }
    
    var pkgs []Package
    for _, line := range strings.Split(out, "\n") {
        if strings.HasPrefix(line, "package:") {
            name := strings.TrimPrefix(line, "package:")
            pkgs = append(pkgs, Package{Name: name})
        }
    }
    return pkgs, nil
}

func getPackageSize(pkg string) int64 {
    out, _ := adb("shell", "pm dump "+pkg, "|", "grep", "versionCode")
    // Simplified — real implementation would parse dumpsys more carefully
    return 0
}

func main() {
    listCmd := flag.NewFlagSet("list", flag.ExitOnError)
    userOnly := listCmd.Bool("user", false, "user-installed only")
    
    searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
    
    exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
    
    if len(os.Args) < 2 {
        fmt.Println("Usage: apm [list|search|export|size]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "list":
        listCmd.Parse(os.Args[2:])
        pkgs, err := listPackages(*userOnly)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        for _, p := range pkgs {
            fmt.Println(p.Name)
        }
        
    case "search":
        searchCmd.Parse(os.Args[2:])
        if searchCmd.NArg() < 1 {
            fmt.Println("Usage: apm search <keyword>")
            os.Exit(1)
        }
        keyword := searchCmd.Arg(0)
        pkgs, _ := listPackages(false)
        for _, p := range pkgs {
            if strings.Contains(strings.ToLower(p.Name), strings.ToLower(keyword)) {
                fmt.Println(p.Name)
            }
        }
        
    case "export":
        exportCmd.Parse(os.Args[2:])
        if exportCmd.NArg() < 1 {
            fmt.Println("Usage: apm export <output.json>")
            os.Exit(1)
        }
        pkgs, _ := listPackages(false)
        data, _ := json.MarshalIndent(pkgs, "", "  ")
        os.WriteFile(exportCmd.Arg(0), data, 0644)
        fmt.Printf("Exported %d packages\n", len(pkgs))
        
    default:
        fmt.Println("Unknown command")
    }
}
