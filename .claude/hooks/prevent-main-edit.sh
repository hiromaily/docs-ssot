#!/bin/sh
# Blocks Edit and Write tool use when the current branch is main or master.

branch=$(git branch --show-current 2>/dev/null)
if [ "$branch" = "main" ] || [ "$branch" = "master" ]; then
  echo "ERROR: Cannot edit files directly on '$branch'. Create a feature branch first." >&2
  exit 1
fi
