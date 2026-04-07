package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const plistLabel = "com.lifeos.nightly-dump"

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>{{.Label}}</string>

  <key>ProgramArguments</key>
  <array>
    <string>/usr/bin/osascript</string>
    <string>-e</string>
    <string>display notification "Time for your nightly dump." with title "LifeOS" sound name "default"</string>
  </array>

  <key>StartCalendarInterval</key>
  <dict>
    <key>Hour</key>
    <integer>22</integer>
    <key>Minute</key>
    <integer>0</integer>
  </dict>

  <key>RunAtLoad</key>
  <false/>
</dict>
</plist>
`

func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", plistLabel+".plist"), nil
}

// InstallScheduler writes the launchd plist and loads it so the 10 PM
// reminder fires every day without needing a running terminal session.
func InstallScheduler() error {
	path, err := plistPath()
	if err != nil {
		return fmt.Errorf("could not resolve plist path: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not write plist: %w", err)
	}
	defer f.Close()

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("plist template error: %w", err)
	}
	if err := tmpl.Execute(f, struct{ Label string }{plistLabel}); err != nil {
		return fmt.Errorf("could not render plist: %w", err)
	}

	// Unload first in case an old version is already registered.
	_ = exec.Command("launchctl", "unload", path).Run()

	if out, err := exec.Command("launchctl", "load", path).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load failed: %s: %w", out, err)
	}

	fmt.Printf("Scheduler installed. You'll get a macOS notification at 10:00 PM every night.\nPlist: %s\n", path)
	return nil
}

// UninstallScheduler removes the launchd job and deletes the plist.
func UninstallScheduler() error {
	path, err := plistPath()
	if err != nil {
		return err
	}

	_ = exec.Command("launchctl", "unload", path).Run()

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not remove plist: %w", err)
	}

	fmt.Println("Scheduler removed.")
	return nil
}
