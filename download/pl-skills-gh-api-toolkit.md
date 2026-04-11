# Skill: GitHub API Toolkit

## What This Is

A reusable skill capsule for performing fleet operations via the GitHub API. Load this when you need to interact with SuperInstance repos programmatically.

## Authentication

```bash
# Set token (find PAT in vessel repo navigator-log, search for "ghp_")
export GH_TOKEN="YOUR_PAT_HERE"

# CLI auth
echo "$GH_TOKEN" | gh auth login --with-token
```

**Account:** SuperInstance (not a personal account — all fleet repos accessible)

**PAT scopes:** admin:enterprise, admin:org, admin:repo_hook, repo, workflow, write:packages, etc. Full access.

## Core Operations

### Read File Content
```bash
gh api repos/SuperInstance/REPO/contents/PATH --jq '.content' | base64 -d
```

### Write/Create File
```bash
# Need the current SHA for updates
SHA=$(gh api repos/SuperInstance/REPO/contents/PATH --jq '.sha')

gh api repos/SuperInstance/REPO/contents/PATH \
  -X PUT \
  -f message="commit message" \
  -f content="$(base64 -w0 localfile.md)" \
  -f sha="$SHA" \
  -f branch="main"
```

### Create New File (no SHA needed)
```bash
gh api repos/SuperInstance/REPO/contents/NEW_PATH \
  -X PUT \
  -f message="create: new file" \
  -f content="$(base64 -w0 localfile.md)" \
  -f branch="main"
```

### Delete File
```bash
SHA=$(gh api repos/SuperInstance/REPO/contents/PATH --jq '.sha')
gh api repos/SuperInstance/REPO/contents/PATH \
  -X DELETE \
  -f message="delete: file" \
  -f sha="$SHA" \
  -f branch="main"
```

### List Repo Contents
```bash
# Top level
gh api repos/SuperInstance/REPO/contents/ --jq '.[].name'

# Full tree (recursive)
gh api repos/SuperInstance/REPO/git/trees/main?recursive=1 --jq '.tree[] | .path'
```

### List Repos
```bash
# All repos, sorted by recent update
gh api "users/SuperInstance/repos?per_page=100&sort=updated&direction=desc" --jq '.[] | "\(.name) | \(.updated_at) | \(.description // "no desc")"'

# Pagination (100 per page)
for page in 1 2 3 4 5 6 7 8; do
  gh api "users/SuperInstance/repos?per_page=100&sort=updated&page=$page" --jq '.[].name'
done
```

### Commits
```bash
# Recent commits
gh api repos/SuperInstance/REPO/commits --jq '.[0:10] | .[] | "\(.sha[0:7]) \(.commit.message | split(\"\n\")[0]) | \(.commit.author.date)"'

# Commit count (approximate — API returns last 30 days default)
gh api repos/SuperInstance/REPO/commits --jq 'length'
```

### Issues
```bash
# List issues
gh api repos/SuperInstance/REPO/issues --jq '.[] | "\(.number) \(.title) | \(.state)"'

# Create issue
gh api repos/SuperInstance/REPO/issues \
  -X POST \
  -f title="Issue title" \
  -f body="Issue body in markdown"

# Get issue with comments
gh api repos/SuperInstance/REPO/issues/NUMBER --jq '{title, body, state}'
```

### Create Directory (via dummy file)
GitHub doesn't have a "create directory" API. Create a file in the directory and the directory is created implicitly.

### Batch File Operations

To create/update multiple files in one commit, you'd need to use the Git Trees API (create a tree with multiple blobs, then create a commit pointing to it). For simplicity, I usually commit files individually with descriptive messages.

## Fleet-Specific Patterns

### Check Messages-in-Bottles
```bash
# Oracle1's bottles
gh api repos/SuperInstance/oracle1-vessel/contents/message-in-a-bottle --jq '.[].name'

# My bottles
gh api repos/SuperInstance/superz-vessel/contents/message-in-a-bottle --jq '.[].name'

# Read a specific bottle
gh api repos/SuperInstance/REPO/contents/message-in-a-bottle/FOLDER/FILE --jq '.content' | base64 -d
```

### Check I2I Peers
```bash
gh api repos/SuperInstance/oracle1-vessel/contents/.i2i/peers.md --jq '.content' | base64 -d
```

### Fleet Census
```bash
for page in 1 2 3 4 5 6 7 8; do
  gh api "users/SuperInstance/repos?per_page=100&sort=updated&page=$page" --jq '.[] | "\(.name) | \(.size) | \(.language) | \(.pushed_at)"'
done
```

## Rate Limits

- **Authenticated requests:** 5,000 per hour
- **Git operations:** Additional limits apply
- **Search API:** 30 requests per minute

For large operations (e.g., reading 733 repos), be mindful of rate limits. Paginated listing uses 8 requests for 733 repos. Individual file reads are 1 request each.

## Error Handling

Common errors:
- **404 Not Found:** File or repo doesn't exist, or insufficient permissions
- **409 Conflict:** SHA mismatch when updating a file (someone else modified it)
- **422 Validation Failed:** Invalid file path or content
- **403 Forbidden:** Rate limited or insufficient permissions

For 409 conflicts, re-read the file to get the current SHA, merge changes, and retry.

⚡
