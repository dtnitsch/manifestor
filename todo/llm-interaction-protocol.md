# LLM Interaction Protocol

**Status:** Planned for v0.5+
**Priority:** P0 (Core value proposition)
**Created:** 2026-01-06

---

## Problem Statement

**Current state:**
- Users paste entire manifests into LLMs (wasteful)
- No standard protocol for LLM-to-manifest interaction
- LLMs don't know how to lazy-load efficiently
- No guidance on token optimization strategies

**Desired state:**
- Single-file entry point for LLMs
- Self-documenting protocol
- Automatic detection of available files
- Built-in token savings tracking
- Deterministic lazy-loading instructions

---

## The Vision: `manifestor.llm.yaml`

### User Workflow (CLI-based LLMs)

**Human to LLM:**
```
Load /path/to/project/manifestor.llm.yaml and answer:
"Which directory contains the authentication logic?"
```

**LLM behavior:**
1. Reads `manifestor.llm.yaml` (300 tokens)
2. Sees protocol instructions + folder structure
3. Identifies relevant directory (e.g., "auth/")
4. Follows lazy-load instructions to read `manifest/auth.yaml`
5. Answers question (total: ~2,000 tokens vs 46,000)

---

## File Structure: `manifestor.llm.yaml`

### Concept: Self-Documenting Protocol File

```yaml
# manifestor.llm.yaml - LLM Interaction Protocol
# This file is designed to be read by LLMs first
# It provides structure + instructions for efficient exploration

protocol:
  version: "1.0"
  instructions: |
    You are reading a manifestor-generated repository index.

    LAZY LOADING PROTOCOL:
    1. This file contains the folder structure (minimal tokens)
    2. For details on a specific folder, load: manifest/{folder}.yaml
    3. For full repository scan, load: manifest.yaml
    4. Use yq for filtering (examples below)

    TOKEN OPTIMIZATION:
    - This file: ~300 tokens
    - Single folder: ~1,500 tokens avg
    - Full manifest: ~46,000 tokens
    - Always start here, load selectively

    INVARIANTS:
    - All paths are relative to root
    - Folder stats are rollups (include subdirectories)
    - Files are sorted by path
    - Split files may not exist (check before loading)

  available_files:
    - path: manifestor.llm.yaml
      description: "This file - start here"
      tokens: 300

    - path: manifest-folders.yaml
      description: "Folder structure only (alias of this file)"
      tokens: 300

    - path: manifest.yaml
      description: "Full manifest - all files and directories"
      tokens: 46305
      use_when: "Global queries, full scans, cross-directory analysis"

    - path: manifest/
      description: "Split files by top-level directory (if enabled)"
      tokens_per_file: ~1500
      use_when: "Focused investigation of specific directories"
      check_exists: true

  yq_examples:
    structure_browsing: |
      # List all directories with file counts
      yq '.folders[] | {path, files}' manifestor.llm.yaml

    find_directory: |
      # Find directories matching pattern
      yq '.folders[] | select(.path | test("auth|login|user"))' manifestor.llm.yaml

    load_specific: |
      # After identifying target directory from this file:
      yq '.' manifest/utilities.yaml

    filter_extensions: |
      # Find Python files in a specific directory
      yq '.nodes[] | select(.path | test("\\.py$"))' manifest/utilities.yaml

    global_search: |
      # Only use manifest.yaml for queries that need everything
      yq '.nodes[] | select(.path | test("config"))' manifest.yaml

root: /Users/daniel.nitsch/Work/titan-mt
generated_at: 2026-01-06T15:11:20Z
manifestor_version: "0.5.0"

token_savings:
  baseline: 46305  # Full manifest in one read
  optimized_workflow: 1560  # Structure + one directory
  savings_percentage: 96.6%
  cost_savings_usd: 0.134  # Per query @ $3/M tokens

folders:
  - path: .
    files: 47
    subfolders: 16
    size_bytes: 88160
    top_extensions: [.py, .sh, .txt]

  - path: utilities
    files: 63
    subfolders: 33
    size_bytes: 150000
    top_extensions: [.py, .yaml]
    split_file: manifest/utilities.yaml  # Hint: split available

  - path: yggdrasil
    files: 16
    subfolders: 33
    size_bytes: 85000
    top_extensions: [.py, .conf]
    split_file: manifest/yggdrasil.yaml

  - path: pantheon
    files: 40
    subfolders: 4
    size_bytes: 120000
    top_extensions: [.py, .yaml]
    split_file: manifest/pantheon.yaml

  # ... small directories (no split_file field)

  - path: bin
    files: 7
    subfolders: 0
    size_bytes: 15000
    top_extensions: [.sh]
    # No split_file - too small, included in manifest.yaml only
```

---

## Key Features

### 1. Self-Documenting Protocol

**Instructions embedded in the file itself:**
- LLM reads once, understands protocol
- No external docs needed
- Version-stamped (future-proof)

### 2. File Discovery

**`available_files` section tells LLM:**
- What files exist
- Token cost of each
- When to use each
- Whether to check existence first

### 3. Token Tracking

**Built-in cost awareness:**
- Shows baseline (full manifest)
- Shows optimized workflow
- Calculates savings percentage
- Provides USD cost estimates

### 4. Lazy Load Hints

**`split_file` field in folders:**
- Tells LLM which directories have split files
- Absence means "use manifest.yaml"
- Deterministic - no guessing

### 5. yq Examples

**Common patterns embedded:**
- LLM can copy/adapt examples
- Reduces query errors
- Teaches best practices

---

## LLM Interaction Examples

### Example 1: Structure Exploration

**User:**
```
Load manifestor.llm.yaml and tell me about the repository structure.
```

**LLM:**
```
I've loaded the repository structure (300 tokens).

This is a Python monorepo with 16 top-level directories:
- utilities/ (63 files, mostly .py and .yaml)
- yggdrasil/ (16 files, .py and .conf)
- pantheon/ (40 files, Python codebase)
...

The largest directories are utilities, pantheon, and yggdrasil.

Total: 47 root files + 129 subdirectories

Would you like details on any specific directory?
```

**Token cost:** 300 (vs 46,000 for full manifest)

---

### Example 2: Targeted Investigation

**User:**
```
Load manifestor.llm.yaml.
Which directory likely contains API endpoints?
```

**LLM:**
```
Based on the folder structure, I see:
- endpoint/ (10 files, 10 subdirs)
- yggdrasil/ (16 files - might be API-related)

Let me load the endpoint directory for details.
[Loads manifest/endpoint.yaml - 273 lines, ~820 tokens]

The endpoint/ directory contains:
- endpoint/rest/ (REST API handlers)
- endpoint/graphql/ (GraphQL resolvers)
...
```

**Token cost:** 300 (structure) + 820 (endpoint) = 1,120 (97% savings)

---

### Example 3: Cross-Directory Query

**User:**
```
Load manifestor.llm.yaml.
Find all Python test files across the entire repository.
```

**LLM:**
```
This requires a global search across all directories.
The protocol indicates I should use manifest.yaml for this.

[Loads manifest.yaml - 46,305 tokens]

Found 127 test files:
- utilities/tests/*.py (45 files)
- yggdrasil/tests/*.py (22 files)
...
```

**Token cost:** 46,305 (appropriate for global query)

---

## Workflow Detection Logic

**LLM decision tree (embedded in protocol):**

```
1. Read manifestor.llm.yaml (always)
2. Determine query scope:

   a) Structure/navigation query?
      → Use folder data in manifestor.llm.yaml (300 tokens)

   b) Focused on 1-2 directories?
      → Check for split_file field
      → Load manifest/{dir}.yaml (~1,500 tokens each)

   c) Cross-directory or global query?
      → Load manifest.yaml (46,000 tokens)

   d) Unsure?
      → Start with manifestor.llm.yaml
      → Ask user for clarification
```

---

## ChatGPT / Web LLM Adaptation (v0.6)

**Challenge:** No file system access, copy/paste only

**Solution: All-in-one paste format**

```yaml
# manifestor.llm.yaml (for ChatGPT paste)

protocol:
  instructions: |
    This is a condensed manifest for copy/paste to web LLMs.

    LAZY LOADING (manual):
    1. This shows folder structure
    2. To see a specific folder: ask user to paste manifest/{folder}.yaml
    3. For full details: ask user to paste manifest.yaml

    USER INSTRUCTIONS:
    Run: manifestor --format llm
    This generates manifestor.llm.yaml (paste it first)
    Then paste additional files as requested by the LLM

# ... rest of structure ...

paste_instructions_for_user: |
  If the LLM asks for more details:

  For a specific directory:
  $ cat manifest/utilities.yaml | pbcopy
  (then paste into chat)

  For full manifest:
  $ cat manifest.yaml | pbcopy
  (then paste into chat)
```

---

## Implementation Checklist

### v0.5.0: Core Protocol

- [ ] Create `manifestor.llm.yaml` generator
- [ ] Add protocol version field
- [ ] Embed lazy-load instructions
- [ ] Add token cost tracking
- [ ] Include yq examples
- [ ] Add `split_file` hints
- [ ] Calculate token savings vs baseline

### v0.5.1: Enhanced Discovery

- [ ] Auto-detect available split files
- [ ] Generate `available_files` section dynamically
- [ ] Add confidence estimates (which file to use)
- [ ] Include common query patterns for the repo type

### v0.6.0: ChatGPT Support

- [ ] Generate paste-friendly format
- [ ] Add user instructions for manual lazy-load
- [ ] Create condensed version (<2000 tokens)
- [ ] Test with Claude/ChatGPT web interfaces

---

## Token Savings Tracking

**Built into every manifestor.llm.yaml:**

```yaml
token_savings:
  baseline: 46305           # Full manifest, single read

  workflows:
    structure_only: 300     # Just folders
    single_directory: 1500  # Folders + one split
    two_directories: 3000   # Folders + two splits
    full_scan: 46305        # Everything

  estimated_savings:
    typical_query: 96%      # Most queries are focused
    cost_per_query: $0.004  # vs $0.139 for full manifest
    annual_savings: $49.30  # Assumes 100 queries/month

  measured_performance:
    generation_time: 0.123s
    file_count: 7
    total_size: 162KB
    split_enabled: true
```

**Why track this?**
- Users see value immediately
- LLMs learn token-efficient patterns
- Justifies the splitting complexity

---

## Config Options

```yaml
# manifestor-config.yaml

output:
  llm:
    enabled: true  # Generate manifestor.llm.yaml
    include_instructions: true
    include_yq_examples: true
    include_token_estimates: true

    # Condensed mode for ChatGPT paste
    condensed: false
    max_tokens: 2000  # For condensed mode
```

---

## Success Criteria

**v0.5 is successful if:**

1. LLMs can load manifestor.llm.yaml and understand protocol (no external docs)
2. Token savings are tracked and visible
3. Lazy-loading happens naturally (LLMs choose right files)
4. Users report <5% of queries need full manifest.yaml
5. Documentation is embedded (file is self-explanatory)

---

## Future: LLM Tool Integration

**v0.7+: Native LLM tool support**

```json
{
  "name": "load_manifest_folder",
  "description": "Load details for a specific repository folder",
  "parameters": {
    "folder": "utilities",
    "manifest_path": "/path/to/manifestor.llm.yaml"
  }
}
```

LLMs with function calling could:
- Read manifestor.llm.yaml first
- Call `load_manifest_folder("utilities")` automatically
- Never load full manifest unless needed

**But that's phase 2.** Start with self-documenting files.

---

## Related Work

- `todo/output-format-yaml.md` - YAML as default
- `todo/manifest-optimization-strategy.md` - Token reduction
- `docs/examples.md` - yq query patterns

---

## Open Questions

1. **Should manifestor.llm.yaml alias manifest-folders.yaml?**
   - Pro: One less file
   - Con: Name clarity ("folders" vs "llm")
   - **Recommendation:** Make llm.yaml primary, folders.yaml alias

2. **How much protocol instruction is too much?**
   - Current: ~200 tokens of instructions
   - Risk: Bloat the minimal file
   - **Recommendation:** Keep ≤20% of total file

3. **Should we version the protocol separately from manifestor?**
   - `protocol.version: "1.0"` vs `manifestor.version: "0.5.0"`
   - **Recommendation:** Yes, allows protocol evolution

4. **ChatGPT condensed mode - how aggressive?**
   - Option A: Top 10 folders only
   - Option B: Smart sampling by size
   - **Recommendation:** Configurable threshold

---

## Implementation Notes

**File generation order:**
1. Scan filesystem (existing)
2. Build manifest (existing)
3. Write manifest.yaml (existing)
4. **NEW:** Analyze manifest in-memory
5. **NEW:** Generate manifestor.llm.yaml with protocol + stats
6. **NEW:** Calculate token savings
7. **NEW:** Write split files (if enabled)

**Code location:**
```
internal/
  llm/
    protocol.go      # LLM protocol file generator
    instructions.go  # Embedded instruction templates
    tokenizer.go     # Token estimation
    savings.go       # Cost calculation
```

---

## The Big Idea

**manifestor.llm.yaml is not just a file.**

It's a **protocol specification** + **data** + **instructions** in one.

LLMs read it and know:
- What they're looking at
- How to use it efficiently
- Where to find more details
- How much they're saving

**Self-documenting, self-optimizing, deterministic.**

This is the killer feature that makes manifestor the standard for LLM-filesystem interaction.
