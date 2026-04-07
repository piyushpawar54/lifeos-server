# nightly-dump-app

A CLI tool for the **Mind** department of LifeOS.

Before sleep, open this app and brain-dump everything on your mind — thoughts, worries, ideas, tasks, reflections. No structure required. Just get it out of your head and into a file.

## What it does

- Prompts you to free-write anything on your mind
- Saves the dump as a timestamped `.txt` file in a local `dumps/` directory
- Keeps a running archive of nightly entries over time

## Why

A cluttered mind makes for poor sleep and poor decisions. The nightly dump is a ritual to offload mental noise before rest, creating a record you can revisit during weekly reviews.

## Usage

```bash
go run main.go
```

Type freely. Press Enter twice to save and exit.

## Output

Dumps are saved to `dumps/dump_YYYY-MM-DD_HH-MM-SS.txt`.

## Part of LifeOS

This app is one module in the [LifeOS](https://github.com/piyushpawar54/lifeos-server) personal operating system.
