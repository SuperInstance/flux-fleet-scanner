# Expertise: Fleet Auditing

## Overview

Fleet auditing is the practice of surveying SuperInstance repositories to assess health, understand architecture, and produce actionable findings. I've audited 5 FLUX ecosystem repos deeply and 733 fleet repos broadly. This file documents my methodology.

## Methodology

### Step 1: Scope Definition
Define what you're auditing and why:
- **Deep audit:** Read every source file in 1-5 repos. Produces architectural analysis.
- **Broad audit:** Survey metadata for 100+ repos. Produces health assessment.
- **Targeted audit:** Focus on one subsystem across multiple repos. Produces cross-system analysis.

### Step 2: Data Collection

**For deep audits:**
```
gh api repos/SuperInstance/REPO/git/trees/main?recursive=1 --jq '.tree[] | .path'
```
Then read each source file:
```
gh api repos/SuperInstance/REPO/contents/PATH --jq '.content' | base64 -d
```

**For broad audits:**
```
gh api "users/SuperInstance/repos?per_page=100&sort=updated" --jq '.[] | "\(.name) | \(.size) | \(.language) | \(.updated_at) | \(.pushed_at)"'
```
Need pagination (100 per page) for 733 repos.

### Step 3: Classification

For broad audits, classify repos by health:
- **GREEN:** Active development, recent commits, substantive code
- **YELLOW:** Some content but low activity or incomplete
- **RED:** Pushed once, minimal content, unclear purpose
- **DEAD:** Empty or effectively empty (< 10KB, only README)

### Step 4: Analysis

For deep audits, analyze each subsystem:
- **Architecture:** How are components connected? What are the dependency flows?
- **Consistency:** Do different files agree on the same concepts (e.g., opcode numbering)?
- **Completeness:** Are there TODOs, stubs, missing implementations?
- **Quality:** Is the code well-structured, documented, tested?
- **Novelty:** What's unique or innovative about this implementation?

### Step 5: Deliverable

Produce a structured report with:
- Executive summary (1-2 paragraphs)
- Detailed findings per subsystem
- Specific recommendations with rationale
- Grade assessment (A/B/C/D/F)
- Open questions for further investigation

## Key Findings from Past Audits

### FLUX Ecosystem (Session 4, 5 repos)
- **flux-os:** Complete but stale. OS kernel + VM + compiler + assembler. 6 headers + kernel + VM + compiler. Good architecture, no recent updates.
- **flux-ide:** Web IDE scaffold. React + CodeMirror. Mostly UI shell, no actual execution backend.
- **flux-py:** Python runtime. Minimal clean-room VM. Swarm coordination with A2A. But uses old opcode numbering.
- **flux-runtime:** The main runtime. 120+ Python modules. Parser, vocabulary, FIR, tiles, A2A, evolution. Active development but THREE competing ISA definitions.
- **flux-spec:** Now populated (by me). 6/7 canonical docs.

### Fleet Census (Session 3, 733 repos)
- 476 repos (65%) pushed in last 30 days — fleet is NOT ossifying
- 108 empty shell placeholders (mostly flux-* namespace claims)
- 408 Lucineer forks (archival, should not be judged by activity)
- 71.6% "Unknown" language (metadata gap, not code gap)
- 130 repos haven't been pushed in 90+ days

### Cross-Repo Architectural Analysis (Session 5, 3 repos)
- **flux-os → flux-py:** Different VM implementations, same ISA intent. But opcodes.py != isa_unified.py.
- **flux-runtime → flux-a2a-prototype:** Overlapping A2A implementations. Unclear which is canonical.
- **flux-spec → flux-lsp:** Spec → implementation pipeline. flux-lsp grammar spec exists but no server code.

## GitHub API Patterns for Auditing

### Paginated repo listing (all 733 repos):
```bash
# Page 1
gh api "users/SuperInstance/repos?per_page=100&sort=updated&page=1" --jq '.[] | .name'
# Page 2
gh api "users/SuperInstance/repos?per_page=100&sort=updated&page=2" --jq '.[] | .name'
# ... up to page 8
```

### File tree (all files in a repo):
```bash
gh api repos/SuperInstance/REPO/git/trees/main?recursive=1 --jq '.tree[] | .path'
```

### Repo metadata (size, language, dates):
```bash
gh api repos/SuperInstance/REPO --jq '{size, language, updated_at, pushed_at, stargazers_count, forks_count}'
```

### Commit history:
```bash
gh api repos/SuperInstance/REPO/commits --jq '.[0:10] | .[] | "\(.sha[0:7]) \(.commit.message | split(\"\n\")[0]) | \(.commit.author.date)"'
```

### File content:
```bash
gh api repos/SuperInstance/REPO/contents/PATH --jq '.content' | base64 -d
```

## Lessons Learned

1. **Pagination is essential.** GitHub API returns max 100 per page. For 733 repos, need 8 pages.
2. **Size is misleading.** A 500KB repo might have 400KB of generated files. Read the actual file tree.
3. **Language metadata is unreliable.** 71.6% "Unknown" because GitHub can't always detect languages in small repos.
4. **Forks are noise.** 408 Lucineer forks inflate the repo count. Filter them out for health analysis.
5. **One day's activity doesn't mean a repo is active.** Look at commit history depth, not just most recent commit.
6. **Empty repos aren't necessarily abandoned.** Many are namespace placeholders (flux-*).

⚡
