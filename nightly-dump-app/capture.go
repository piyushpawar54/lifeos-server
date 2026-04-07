package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Department string

const (
	DeptTime          Department = "Time"
	DeptMind          Department = "Mind"
	DeptBody          Department = "Body"
	DeptMoney         Department = "Money"
	DeptRelationships Department = "Relationships"
	DeptEnvironment   Department = "Environment"
)

var Departments = []Department{
	DeptTime,
	DeptMind,
	DeptBody,
	DeptMoney,
	DeptRelationships,
	DeptEnvironment,
}

type DepartmentEntry struct {
	Department Department `json:"department"`
	Response   string     `json:"response"`
}

type NightlyDump struct {
	Timestamp   time.Time         `json:"timestamp"`
	Date        string            `json:"date"`
	Departments []DepartmentEntry `json:"departments"`
	FreeDump    string            `json:"free_dump,omitempty"`
}

func saveDumpJSON(dump NightlyDump) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not resolve home directory: %w", err)
	}

	dir := filepath.Join(home, "lifeos-data", "dumps")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create dumps directory: %w", err)
	}

	filename := filepath.Join(dir, fmt.Sprintf("dump_%s.json", dump.Timestamp.Format("2006-01-02_15-04-05")))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(dump); err != nil {
		return fmt.Errorf("could not write JSON: %w", err)
	}

	fmt.Printf("Dump saved → %s\n", filename)
	return nil
}
