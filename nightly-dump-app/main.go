package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type DumpEntry struct {
	Timestamp time.Time
	Content   string
}

func main() {
	fmt.Println("=== LifeOS Nightly Dump ===")
	fmt.Println("Brain dump everything on your mind. Press Enter twice to finish.")
	fmt.Println()

	entry := collectDump()

	if strings.TrimSpace(entry.Content) == "" {
		fmt.Println("Nothing captured. Goodnight.")
		return
	}

	if err := saveDump(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving dump: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDump saved at %s\n", entry.Timestamp.Format("2006-01-02 15:04"))
}

func collectDump() DumpEntry {
	scanner := bufio.NewScanner(os.Stdin)
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

	return DumpEntry{
		Timestamp: time.Now(),
		Content:   strings.Join(lines, "\n"),
	}
}

func saveDump(entry DumpEntry) error {
	dir := "dumps"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/dump_%s.txt", dir, entry.Timestamp.Format("2006-01-02_15-04-05"))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Date: %s\n\n", entry.Timestamp.Format("Monday, January 2 2006 — 15:04"))
	fmt.Fprintf(f, "%s\n", entry.Content)
	return nil
}
