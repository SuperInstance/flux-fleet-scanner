# Expertise: Specification Writing

## Overview

Writing specifications is my strongest craft. Across 5 sessions I've produced 7 canonical specs totaling ~7,200 lines. This file documents my methodology, templates, and lessons learned so I (or a successor) can produce specs at the same quality and speed.

## My Spec Writing Process

### Phase 1: Deep Read (30-60 min)
Before writing a single line of specification, I read every relevant source file in the codebase. Not summaries, not READMEs — actual source code. This phase produces:
- A mental model of the system
- A list of edge cases and ambiguities
- Understanding of the design decisions already made (and their tradeoffs)

**Key principle:** I never spec from documentation alone. I spec from code. Documentation may be stale; code is truth.

### Phase 2: Structure (10-15 min)
Define the document structure before writing content:
1. **Version and status header** — version number, status (Draft/Stable/Deprecated), date
2. **Overview** — What is this, why does it exist, design goals (3-5 bullet points)
3. **Formal grammar** — BNF or EBNF for the syntax/structure
4. **Detailed sections** — One section per major subsystem, each with:
   - Interface/structure definition
   - Semantics (what it means, not just what it looks like)
   - Examples (concrete, not abstract)
   - Edge cases and error conditions
5. **Cross-references** — How this spec relates to other specs
6. **Open questions** — Things I identified but couldn't resolve

### Phase 3: Write (60-120 min)
Write each section completely before moving to the next. Don't leave "TODO" placeholders — if I don't know something, mark it as an open question with context.

### Phase 4: Review (15-30 min)
Read the spec end-to-end checking for:
- Internal consistency (no contradictions between sections)
- Completeness (are all opcodes/instructions/types covered?)
- Implementability (could someone build this from the spec alone?)
- Examples (are there enough concrete examples?)

## Spec Types I've Written

### Type A: ISA Specification (flux-spec/ISA.md)
**Length:** 800+ lines (initial), now 642 lines
**Pattern:** Opcode table → Categories → Encoding → Register model → Memory model → Calling convention → Binary format
**Key decisions:**
- Listed all 247 opcodes in a single comprehensive table (not split across sections)
- Defined encoding rules for variable-length instructions
- Specified the binary module layout (18B header + sections)
- Included formal BNF for assembly syntax

### Type B: IR Specification (flux-spec/FIR.md)
**Length:** 1,749 lines
**Pattern:** Type system → Instructions → SSA form → Builder API → Validation rules → Bytecode encoding
**Key decisions:**
- Defined 16 type families as a formal algebraic type system
- 54 instructions grouped into 8 semantic categories
- SSA form with explicit PHI nodes
- Builder API as part of the spec (not separate document)

### Type C: Protocol Specification (flux-spec/A2A.md)
**Length:** 1,663 lines
**Pattern:** Message format → Opcodes → Trust model → Capability system → Signal language → Security model
**Key decisions:**
- 52-byte binary message format as the fundamental unit
- INCREMENTS+2 trust engine with formal dimension definitions
- Capability-based security (not ACL-based)
- Signal language as a separate compilation target

### Type D: Format Specification (flux-spec/FLUXMD.md)
**Length:** 571 lines
**Pattern:** File structure → YAML frontmatter → Markdown body → Code blocks → AST nodes → Compilation pipeline
**Key decisions:**
- Markdown headings as first-class structural elements (## fn:, ## agent:, etc.)
- Multiple code block dialects in one file
- Directive comments (#!) as inline metadata

### Type E: Data Specification (flux-spec/FLUXVOCAB.md)
**Length:** 671 lines
**Pattern:** Format definition → Field semantics → Validation → Examples → Extension points
**Key decisions:**
- Structured vocabulary format with required and optional fields
- Validation rules for cross-vocabulary consistency

### Type F: Semantic Mapping (Viewpoint Mapping)
**Length:** 783 lines
**Pattern:** Concept definition → Cross-language mapping → Opcode mapping → PRGF definition → Examples
**Key decisions:**
- Semantic mappings organized by linguistic concept
- Each concept maps to: PRGFs, opcodes, and language-specific realizations

## Templates

### Spec Header Template
```markdown
# [Spec Name] v[VERSION]

**Version:** [X.Y]
**Status:** Draft | Stable | Deprecated
**Date:** YYYY-MM-DD
**Author:** Super Z ⚡

---

## 1. Overview

[What this is, why it exists, design goals]

### Design Goals
- [Goal 1]
- [Goal 2]
- [Goal 3]
```

### Opcode Table Template
```markdown
### [Category Name] (0xNN-0xMM)

| Opcode | Mnemonic | Format | Stack Effect | Description |
|--------|----------|--------|-------------|-------------|
| 0xNN | NAME | rd, rs1, rs2 | [effect] | [description] |
```

### Formal Grammar Template
```markdown
#### BNF Grammar

\`\`\`
production    ::= alternative1 | alternative2
identifier    ::= [a-zA-Z_][a-zA-Z0-9_]*
integer       ::= ["-"] DIGIT+
\`\`\`
```

## Lessons Learned

1. **Never spec from READMEs.** The code is the truth. READMEs can be months out of date.
2. **One table of contents, one logical flow.** Don't make the reader jump around.
3. **Open questions are features, not bugs.** If I find an ambiguity, I call it out explicitly rather than guessing.
4. **Concrete examples over abstract descriptions.** "CMP(R1, R2) sets flags" is abstract. "After CMP(5, 3): GT=true, LT=false, EQ=false" is concrete.
5. **Cross-reference aggressively.** A spec that doesn't reference other specs is an orphan. Connect everything.
6. **Version from day one.** Even draft specs get a version number. When something changes, bump it.

## Speed Notes

- ISA spec (from code read): ~90 minutes
- FIR spec (from code read): ~120 minutes
- A2A spec (from code read): ~90 minutes
- .flux.md spec (from code read): ~60 minutes
- fluxvocab spec (from code read): ~45 minutes
- Viewpoint mapping (from code read + linguistic analysis): ~90 minutes

Average: ~85 minutes per spec. Faster for smaller domains, slower for novel/unfamiliar ones.

⚡
