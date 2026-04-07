package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	install := flag.Bool("install", false, "Install the 10 PM daily reminder on macOS")
	uninstall := flag.Bool("uninstall", false, "Remove the 10 PM daily reminder")
	flag.Parse()

	if *install {
		if err := InstallScheduler(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *uninstall {
		if err := UninstallScheduler(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	runDump()
}

func runDump() {
	now := time.Now()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== LifeOS Nightly Dump ===")
	fmt.Printf("%s\n\n", now.Format("Monday, January 2 2006 — 15:04"))
	fmt.Println("Work through each department. Press Enter twice to move on.")
	fmt.Println()

	var entries []DepartmentEntry

	questions := map[Department]string{
		DeptTime:          "Time — Did today go as planned? Where did your time go? What must happen tomorrow?",
		DeptMind:          "Mind — What's taking up mental space? Open loops, anxieties, ideas you don't want to lose?",
		DeptBody:          "Body — How was your energy today? Sleep, movement, food — anything to note?",
		DeptMoney:         "Money — Any spending or financial decisions today? Anything you're sitting on?",
		DeptRelationships: "Relationships — Who did you connect with? Anyone you've been meaning to reach out to?",
		DeptEnvironment:   "Environment — How did your space feel? Anything to fix, clear, or organise?",
	}

	for _, dept := range Departments {
		fmt.Printf("── %s\n", questions[dept])
		response := collectBlock(scanner)
		entries = append(entries, DepartmentEntry{
			Department: dept,
			Response:   response,
		})
		fmt.Println()
	}

	fmt.Println("── Anything else on your mind?")
	freeDump := collectBlock(scanner)

	dump := NightlyDump{
		Timestamp:   now,
		Date:        now.Format("2006-01-02"),
		Departments: entries,
		FreeDump:    freeDump,
	}

	if err := saveDumpJSON(dump); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving dump: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSaved at %s. Goodnight.\n", now.Format("15:04"))
}

// collectBlock reads lines until two consecutive empty lines are entered.
func collectBlock(scanner *bufio.Scanner) string {
	var lines []string
	emptyCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			emptyCount++
			if emptyCount >= 2 {
				break
			}
		} else {
			emptyCount = 0
		}
		lines = append(lines, line)
	}

	// Trim trailing single blank line that accumulates before the double-enter.
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}
