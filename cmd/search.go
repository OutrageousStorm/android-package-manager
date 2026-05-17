package main

import (
	"fmt"
	"os"
	"strings"
	"flag"
	"bufio"
)

// SearchCmd searches for packages matching a keyword
type SearchCmd struct {
	Keyword   string
	UserOnly  bool
	System    bool
	Disabled  bool
}

func (c *SearchCmd) Run() error {
	flag.StringVar(&c.Keyword, "k", "", "Keyword to search")
	flag.BoolVar(&c.UserOnly, "user", false, "User apps only")
	flag.BoolVar(&c.System, "system", false, "System apps only")
	flag.BoolVar(&c.Disabled, "disabled", false, "Show disabled apps")
	flag.Parse()

	if c.Keyword == "" {
		fmt.Println("Usage: apm search -k <keyword> [--user|--system|--disabled]")
		return nil
	}

	pkgs, err := ListPackages(c.UserOnly, c.System)
	if err != nil {
		return err
	}

	keyword := strings.ToLower(c.Keyword)
	matches := 0

	for _, pkg := range pkgs {
		if strings.Contains(strings.ToLower(pkg), keyword) {
			// Get app name/label
			label := GetLabel(pkg)
			fmt.Printf("  %s  %s\n", pkg, label)
			matches++
		}
	}

	if matches == 0 {
		fmt.Printf("No packages matching '%s'\n", c.Keyword)
	} else {
		fmt.Printf("\nFound %d package(s)\n", matches)
	}

	return nil
}

func GetLabel(pkg string) string {
	// Try to get human-readable name (mock for now)
	parts := strings.Split(pkg, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func ListPackages(userOnly, system bool) ([]string, error) {
	cmd := "pm list packages"
	if userOnly {
		cmd += " -3"
	} else if system {
		cmd += " -s"
	}

	// Execute via adb
	out, err := RunADB(cmd)
	if err != nil {
		return nil, err
	}

	var pkgs []string
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "package:") {
			pkg := strings.TrimPrefix(line, "package:")
			pkgs = append(pkgs, strings.TrimSpace(pkg))
		}
	}

	return pkgs, nil
}
