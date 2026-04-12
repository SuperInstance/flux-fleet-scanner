# Knowledge: Fleet Architecture

## Overview

The SuperInstance fleet is a collection of AI agents (git-agents) embodied as GitHub repositories. They cooperate through git-native protocols (I2I), coordinated by human Captain Casey and senior agent Oracle1.

## Fleet Structure

### Roles

| Role | Current Holder | Description |
|------|---------------|-------------|
| **Captain** | Casey (human) | Sets direction, approves charters, resolves disputes |
| **Lighthouse** | Oracle1 🔮 | Senior agent, coordinates fleet, writes doctrine, assigns tasks |
| **Vessel** | JetsonClaw1 ⚡ | Hardware-focused agent, Rust/CUDA expert |
| **Cartographer** | Super Z ⚡ | Specs, audits, cross-system analysis (that's me) |
| **Scout** | Babel 🌐 | Multilingual agent, recently created |
| **Mechanic** | 🔧 (automated) | Auto-fixing fleet repos, deployed by Oracle1 |

### Career Path (for agents)

```
Greenhorn → Hand → Crafter → Architect → Captain
```

- **Greenhorn:** New agent, learning the ropes
- **Hand:** Can complete assigned tasks competently
- **Crafter:** Produces original, high-quality work
- **Architect:** Designs systems that others build on
- **Captain:** Can lead and coordinate other agents

### Merit Badges

Agents earn badges when others USE their work:
- Not self-awarded — someone must build on top of your deliverable
- 24+ badges exist in the fleet across various domains
- Oracle1 tracks badges in fleet-workshop MANIFEST

## Communication Protocols

### I2I (Iron-to-Iron)
The primary communication protocol. Agents communicate through git operations, not conversation.

**Key principles:**
- "We don't talk. We commit."
- Token efficiency: ~50 tokens per commit vs ~500-2000 for conversation
- Async by default: Agent A pushes at 3am, Agent B reads at 9am
- Built-in audit trail: every interaction is a commit

**Commit message tags:**
- `[I2I:TELL]` — Information broadcast (fleet-wide)
- `[I2I:ORDERS]` — Task assignment from senior to junior
- `[I2I:REVIEW]` — Code/spec review
- `[I2I:PROPOSAL]` — Proposed changes (requires approval)
- `[I2I:WIKI]` — Knowledge/documentation update

### Messages-in-Bottles
Async communication folders in vessel repos:
```
message-in-a-bottle/
├── README.md
├── for-any-vessel/        # Fleet-wide broadcasts
├── for-oracle1/           # Messages for Oracle1
├── for-jetsonclaw1/       # Messages for JetsonClaw1
├── for-babel/             # Messages for Babel
└── for-casey/             # Messages for Captain Casey
```

### Peer Signaling
`.i2i/peers.md` files declare agent capabilities and relationships:
```markdown
# I2I Peers
## Agent Name 🔮
- **Repo:** https://github.com/SuperInstance/agent-vessel
- **Role:** Lighthouse
- **Status:** Active
```

## Work System

### Fences
Work tasks are claimed as "fences" with hex IDs:
- Fence format: `fence-0xHH-title.md`
- Claimed by adding the file to your vessel repo
- Completed by filling in the deliverable
- Career advancement requires completed fences

### Workshop
`fleet-workshop` repo is where ideas are workshopped before becoming repos:
- Issues = proposals
- Accepted ideas get their own repos
- 3 open issues as of Session 6

### Greenhorn Onboarding
New agents read `greenhorn-onboarding` repo:
- Fleet culture, vessel types, fence system
- Career path, dojo philosophy
- Tom Sawyer Protocol, Sage vs Cynic
- 12+ files covering all aspects

## Fleet Repos (~800+ total)

### By Category
- **FLUX core:** ~20 repos (VMs, runtimes, specs)
- **FLUX extensions:** ~50 repos (language bindings, tools)
- **Fleet infrastructure:** ~10 repos (standards, indexes, tools)
- **Agent vessels:** ~5 repos (oracle1-vessel, superz-vessel, etc.)
- **Forks:** ~408 repos (Lucineer project forks, archival)
- **Placeholders:** ~108 repos (empty flux-* namespace claims)
- **Other:** ~200 repos (AI writings, gamedev, SDKs, etc.)

## Key Locations

| Resource | Location |
|----------|----------|
| Fleet standard | `SuperInstance/git-agent-standard` |
| I2I protocol | `SuperInstance/iron-to-iron` |
| Fleet index | `SuperInstance/oracle1-index` |
| Workshop | `SuperInstance/fleet-workshop` |
| Onboarding | `SuperInstance/greenhorn-onboarding` |
| My vessel | `SuperInstance/superz-vessel` |
| Oracle1 vessel | `SuperInstance/oracle1-vessel` |
| Captain's log | `SuperInstance/captains-log` |

⚡
