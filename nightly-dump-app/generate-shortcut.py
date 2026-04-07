#!/usr/bin/env python3
"""
Generates 'LifeOS Nightly Dump.shortcut' as a binary plist.

iOS Shortcut actions:
  1. Get Contents of URL  — fetches dump-prompt.md from GitHub
  2. Copy to Clipboard    — puts the prompt on the clipboard
  3. Show Notification    — instructs user to paste + reminds about save-dump
  4. Open URL             — opens https://claude.ai/new in Safari

To set up the 10 PM daily trigger:
  iOS Shortcuts → Automation → New Automation → Time of Day → 10:00 PM
  → Run Shortcut → LifeOS Nightly Dump → Don't Ask Before Running
"""

import plistlib
import uuid
import os
import sys

PROMPT_URL = (
    "https://raw.githubusercontent.com/piyushpawar54/lifeos-server"
    "/main/nightly-dump-app/dump-prompt.md"
)
CLAUDE_URL = "https://claude.ai/new"


def text_token(s: str) -> dict:
    return {
        "Value": {"string": s},
        "WFSerializationType": "WFTextTokenString",
    }


def variable_token(output_name: str, output_uuid: str) -> dict:
    """Reference the output of a previous action by UUID."""
    return {
        "Value": {
            "attachmentsByRange": {
                "{0, 1}": {
                    "OutputName": output_name,
                    "OutputUUID": output_uuid,
                    "Type": "ActionOutput",
                }
            },
            "string": "\ufffc",  # object-replacement character placeholder
        },
        "WFSerializationType": "WFTextTokenString",
    }


def new_uuid() -> str:
    return str(uuid.uuid4()).upper()


def build_shortcut() -> dict:
    fetch_uuid = new_uuid()
    copy_uuid  = new_uuid()
    notif_uuid = new_uuid()
    open_uuid  = new_uuid()

    actions = [
        # 1 ── Fetch dump-prompt.md from GitHub
        {
            "WFWorkflowActionIdentifier": "is.workflow.actions.downloadurl",
            "WFWorkflowActionParameters": {
                "UUID": fetch_uuid,
                "WFHTTPMethod": "GET",
                "WFURL": PROMPT_URL,
                "WFShowHeaders": False,
            },
        },
        # 2 ── Copy fetched text to clipboard
        {
            "WFWorkflowActionIdentifier": "is.workflow.actions.setclipboard",
            "WFWorkflowActionParameters": {
                "UUID": copy_uuid,
                "WFInput": variable_token("Contents of URL", fetch_uuid),
            },
        },
        # 3 ── Notify the user
        {
            "WFWorkflowActionIdentifier": "is.workflow.actions.notification",
            "WFWorkflowActionParameters": {
                "UUID": notif_uuid,
                "WFNotificationActionTitle": text_token("LifeOS Nightly Dump"),
                "WFNotificationActionBody": text_token(
                    "Prompt copied to clipboard.\n"
                    "Paste it into Claude.ai and work through your six departments.\n\n"
                    "After the session, run in your terminal:\n"
                    "  cd lifeos-server/nightly-dump-app && go run . save-dump"
                ),
                "WFNotificationActionSound": True,
            },
        },
        # 4 ── Open Claude.ai
        {
            "WFWorkflowActionIdentifier": "is.workflow.actions.openurl",
            "WFWorkflowActionParameters": {
                "UUID": open_uuid,
                "WFInput": text_token(CLAUDE_URL),
            },
        },
    ]

    return {
        "WFWorkflowClientVersion": "2605.0.5",
        "WFWorkflowHasOutputFallback": False,
        "WFWorkflowIcon": {
            "WFWorkflowIconGlyphNumber": 59511,   # moon glyph
            "WFWorkflowIconStartColor": 431817727, # dark blue
        },
        "WFWorkflowImportQuestions": [],
        "WFWorkflowInputContentItemClasses": [],
        "WFWorkflowMinimumClientVersion": 900,
        "WFWorkflowMinimumClientVersionString": "900",
        "WFWorkflowName": "LifeOS Nightly Dump",
        "WFWorkflowNoInputBehavior": {
            "Name": "WFWorkflowNoInputBehaviorAskForInput",
            "Parameters": {},
        },
        "WFWorkflowOutputContentItemClasses": [],
        "WFWorkflowTypes": [],
        "WFWorkflowActions": actions,
    }


def main():
    out = os.path.join(os.path.dirname(os.path.abspath(__file__)),
                       "LifeOS Nightly Dump.shortcut")
    shortcut = build_shortcut()
    with open(out, "wb") as f:
        plistlib.dump(shortcut, f, fmt=plistlib.FMT_BINARY)
    print(f"Written → {out}")


if __name__ == "__main__":
    main()
