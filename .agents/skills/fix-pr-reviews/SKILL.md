---
name: fix-pr-reviews
description: Fix PR review comments for the current branch by reading inline comments, applying fixes, and replying.
---

# Fix PR review comments

# Fix PR review comments

For the current branch, determine the PR number with `gh pr view --json number`.

Fetch inline review comments (diff comments) with:
```
gh api repos/{owner}/{repo}/pulls/{pull_number}/comments
```
These are separate from issue-level comments. Do **not** use `gh pr view --json comments` — that field returns only issue-level (non-diff) comments and will miss inline review comments.

For each comment that has not yet been resolved:

1. Read the referenced file and line to understand the context.
2. Apply the fix directly in the file.
3. Reply to the comment using:
   ```
   gh api repos/{owner}/{repo}/pulls/{pull_number}/comments/{comment_id}/replies -f body="<reply>"
   ```
   Note: the pull number **must** be included in the path — omitting it returns 404.

After all comments are addressed, stage and commit the changes, then push the branch.

Keep replies short and factual: what was changed and why. Do not ask for confirmation before fixing — just fix and reply.
