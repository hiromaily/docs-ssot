#!/bin/sh
# Blocks git commit and git push when the current branch is main or master.

case "$CLAUDE_TOOL_INPUT_COMMAND" in
  *"git commit"* | *"git push"*)
    branch=$(git branch --show-current 2>/dev/null)
    if [ "$branch" = "main" ] || [ "$branch" = "master" ]; then
      echo "ERROR: Direct commit/push to '$branch' is forbidden. Create a feature branch first." >&2
      exit 1
    fi
    ;;
esac
