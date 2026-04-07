# nightly-dump-app

A CLI tool for the **Mind** department of LifeOS.

Every night before sleep, run this app and work through all six LifeOS departments — Time, Mind, Body, Money, Relationships, Environment. Your responses are saved as a timestamped JSON file in `~/lifeos-data/dumps/`. A macOS LaunchAgent fires a notification at 10 PM so you never forget.

---

## Full Flow

```
10:00 PM  macOS notification fires → "Time for your nightly dump."
    ↓
You open terminal and run: go run .
    ↓
App walks you through six departments, one at a time
    ↓
Responses saved to ~/lifeos-data/dumps/dump_YYYY-MM-DD_HH-MM-SS.json
    ↓
Sleep with a clear head
```

---

## Files

| File | Purpose |
|---|---|
| `main.go` | Entry point — orchestrates the six-department interview |
| `capture.go` | Data types and JSON save logic → `~/lifeos-data/dumps/` |
| `scheduler.go` | macOS LaunchAgent installer for the 10 PM daily reminder |
| `dump-prompt.md` | Reference questions for each department |

---

## Setup

**1. Install the daily reminder (once)**

```bash
go run . -install
```

This writes a LaunchAgent plist to `~/Library/LaunchAgents/com.lifeos.nightly-dump.plist` and loads it. You'll receive a macOS notification at 10:00 PM every day.

**2. Run a nightly dump**

```bash
go run .
```

The app will step through each department. Press **Enter twice** to move to the next section.

**3. Remove the reminder**

```bash
go run . -uninstall
```

---

## Output Format

Dumps are saved to `~/lifeos-data/dumps/dump_YYYY-MM-DD_HH-MM-SS.json`.

```json
{
  "timestamp": "2026-04-07T22:03:41Z",
  "date": "2026-04-07",
  "departments": [
    { "department": "Time",          "response": "..." },
    { "department": "Mind",          "response": "..." },
    { "department": "Body",          "response": "..." },
    { "department": "Money",         "response": "..." },
    { "department": "Relationships", "response": "..." },
    { "department": "Environment",   "response": "..." }
  ],
  "free_dump": "..."
}
```

---

## Why JSON

Plain text is human-readable but hard to analyse. JSON lets the future `weekly-review-app` and `insights-app` query your history by department, date range, or keyword without any parsing gymnastics.

---

## Part of LifeOS

This app is one module in the [LifeOS](https://github.com/piyushpawar54/lifeos-server) personal operating system.
