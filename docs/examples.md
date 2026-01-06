# Manifest Query Examples

This guide shows common queries for exploring manifests using `yq` (YAML) and `jq` (JSON).

## Prerequisites

**For YAML manifests:**
```bash
# macOS
brew install yq

# Linux
wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq
chmod +x /usr/bin/yq
```

**For JSON manifests:**
```bash
# macOS
brew install jq

# Most Linux distros include jq in default repos
apt-get install jq  # Debian/Ubuntu
yum install jq      # RHEL/CentOS
```

---

## Basic Queries

### Count Total Nodes

**YAML:**
```bash
yq '.nodes | length' manifest.yaml
```

**JSON:**
```bash
jq '.nodes | length' manifest.json
```

**Example output:**
```
52
```

---

### List All File Paths

**YAML:**
```bash
yq '.nodes[].path' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[].path' manifest.json
```

**Example output:**
```
.
.gitignore
README.md
main.go
...
```

---

## Finding Files and Directories

### Find All Directories

**YAML:**
```bash
yq '.nodes[] | select(.is_dir == true) | .path' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select(.is_dir == true) | .path' manifest.json
```

**Example output:**
```
.
docs
internal
internal/manifest
internal/scanner
```

---

### Find All Files (Non-Directories)

**YAML:**
```bash
yq '.nodes[] | select(.is_dir != true) | .path' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select(.is_dir != true) | .path' manifest.json
```

---

### Find Files by Extension

**YAML:**
```bash
# Find all .go files
yq '.nodes[] | select(.path | test("\\.go$")) | .path' manifest.yaml

# Find all .md files
yq '.nodes[] | select(.path | test("\\.md$")) | .path' manifest.yaml
```

**JSON:**
```bash
# Find all .go files
jq '.nodes[] | select(.path | endswith(".go")) | .path' manifest.json

# Find all .md files
jq '.nodes[] | select(.path | endswith(".md")) | .path' manifest.json
```

---

## Size-Based Queries

### Find Large Files (>1MB)

**YAML:**
```bash
yq '.nodes[] | select(.size_bytes > 1000000) | {path, size: .size_bytes}' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select(.size_bytes > 1000000) | {path, size: .size_bytes}' manifest.json
```

**Example output:**
```yaml
path: manifestor
size: 4059650
```

---

### Find Largest Files (Top 10)

**YAML:**
```bash
yq '.nodes[] | select(.size_bytes) | {path, size: .size_bytes}' manifest.yaml \
  | yq -s 'sort_by(.size) | reverse | .[0:10]'
```

**JSON:**
```bash
jq '.nodes[] | select(.size_bytes) | {path, size: .size_bytes}' manifest.json \
  | jq -s 'sort_by(.size) | reverse | .[0:10]'
```

---

### Find Small Files (<1KB)

**YAML:**
```bash
yq '.nodes[] | select(.size_bytes < 1024) | {path, size: .size_bytes}' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select(.size_bytes < 1024) | {path, size: .size_bytes}' manifest.json
```

---

## Timestamp Queries

### Find Recently Modified Files (Last 7 Days)

**YAML:**
```bash
# Get current Unix timestamp - 7 days
CUTOFF=$(date -v-7d +%s)

yq ".nodes[] | select(.mtime_unix > $CUTOFF) | .path" manifest.yaml
```

**JSON:**
```bash
# Get current Unix timestamp - 7 days
CUTOFF=$(date -v-7d +%s)

jq ".nodes[] | select(.mtime_unix > $CUTOFF) | .path" manifest.json
```

---

### Find Oldest Files

**YAML:**
```bash
yq '.nodes[] | select(.mtime_unix) | {path, mtime: .mtime_unix}' manifest.yaml \
  | yq -s 'sort_by(.mtime) | .[0:10]'
```

**JSON:**
```bash
jq '.nodes[] | select(.mtime_unix) | {path, mtime: .mtime_unix}' manifest.json \
  | jq -s 'sort_by(.mtime) | .[0:10]'
```

---

## Rollup Statistics

### Get Root Directory Statistics

**YAML:**
```bash
yq '.nodes[0].rollup' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[0].rollup' manifest.json
```

**Example output:**
```yaml
total_files: 13
total_descendant_dirs: 10
extensions:
  .go: 12
  .md: 2
  .yaml: 4
size:
  total: 4100207
  min: 28
  max: 4059650
  mean: 315400
  median: 1956
last_modified: 1767663079
```

---

### Find Directories with Most Files

**YAML:**
```bash
yq '.nodes[] | select(.rollup.total_files) | {path, files: .rollup.total_files}' manifest.yaml \
  | yq -s 'sort_by(.files) | reverse | .[0:5]'
```

**JSON:**
```bash
jq '.nodes[] | select(.rollup.total_files) | {path, files: .rollup.total_files}' manifest.json \
  | jq -s 'sort_by(.files) | reverse | .[0:5]'
```

---

### Get Extension Statistics

**YAML:**
```bash
# Show all extensions and their counts from root
yq '.nodes[0].rollup.extensions' manifest.yaml
```

**JSON:**
```bash
# Show all extensions and their counts from root
jq '.nodes[0].rollup.extensions' manifest.json
```

---

### Find Directories by Total Size

**YAML:**
```bash
yq '.nodes[] | select(.rollup.size.total) | {path, size: .rollup.size.total}' manifest.yaml \
  | yq -s 'sort_by(.size) | reverse | .[0:10]'
```

**JSON:**
```bash
jq '.nodes[] | select(.rollup.size.total) | {path, size: .rollup.size.total}' manifest.json \
  | jq -s 'sort_by(.size) | reverse | .[0:10]'
```

---

## Advanced Queries

### Count Files by Type

**YAML:**
```bash
# Count total files by extension across entire manifest
yq '.nodes[].rollup.extensions | to_entries | .[]' manifest.yaml \
  | yq -s 'group_by(.key) | map({ext: .[0].key, count: (map(.value) | add)}) | sort_by(.count) | reverse'
```

**JSON:**
```bash
jq '[.nodes[].rollup.extensions | to_entries | .[]] | group_by(.key) | map({ext: .[0].key, count: (map(.value) | add)}) | sort_by(.count) | reverse' manifest.json
```

---

### Find Directories with Specific Extension

**YAML:**
```bash
# Find all directories containing .go files
yq '.nodes[] | select(.rollup.extensions.".go") | .path' manifest.yaml
```

**JSON:**
```bash
# Find all directories containing .go files
jq '.nodes[] | select(.rollup.extensions[".go"]) | .path' manifest.json
```

---

### Calculate Total Repository Size

**YAML:**
```bash
yq '.nodes[0].rollup.size.total' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[0].rollup.size.total' manifest.json
```

**Convert to human-readable:**
```bash
# YAML
echo "$(yq '.nodes[0].rollup.size.total' manifest.yaml) bytes" | numfmt --to=iec

# JSON
echo "$(jq '.nodes[0].rollup.size.total' manifest.json) bytes" | numfmt --to=iec
```

---

### Find Empty Directories

**YAML:**
```bash
yq '.nodes[] | select(.is_dir == true and .rollup.total_files == 0) | .path' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select(.is_dir == true and .rollup.total_files == 0) | .path' manifest.json
```

---

## Filtering and Combining

### Find Large Go Files

**YAML:**
```bash
yq '.nodes[] | select((.path | test("\\.go$")) and .size_bytes > 10000) | {path, size: .size_bytes}' manifest.yaml
```

**JSON:**
```bash
jq '.nodes[] | select((.path | endswith(".go")) and .size_bytes > 10000) | {path, size: .size_bytes}' manifest.json
```

---

### Get File Count by Directory Depth

**YAML:**
```bash
yq '.nodes[] | select(.is_dir != true) | .path' manifest.yaml \
  | awk -F'/' '{print NF-1}' \
  | sort \
  | uniq -c
```

**JSON:**
```bash
jq -r '.nodes[] | select(.is_dir != true) | .path' manifest.json \
  | awk -F'/' '{print NF-1}' \
  | sort \
  | uniq -c
```

**Example output:**
```
  12 0
   8 1
  15 2
   7 3
```

---

## Export and Processing

### Export to CSV

**YAML:**
```bash
# Export files with size and modification time
echo "path,size_bytes,mtime_unix" > files.csv
yq -r '.nodes[] | select(.size_bytes) | [.path, .size_bytes, .mtime_unix] | @csv' manifest.yaml >> files.csv
```

**JSON:**
```bash
# Export files with size and modification time
echo "path,size_bytes,mtime_unix" > files.csv
jq -r '.nodes[] | select(.size_bytes) | [.path, .size_bytes, .mtime_unix] | @csv' manifest.json >> files.csv
```

---

### Generate Summary Report

**YAML:**
```bash
cat << EOF
Repository Summary
==================
Total nodes: $(yq '.nodes | length' manifest.yaml)
Total files: $(yq '.nodes[0].rollup.total_files' manifest.yaml)
Total size: $(yq '.nodes[0].rollup.size.total' manifest.yaml) bytes
Largest file: $(yq '.nodes[0].rollup.size.max' manifest.yaml) bytes
Generated: $(yq '.generated_at' manifest.yaml)
EOF
```

**JSON:**
```bash
cat << EOF
Repository Summary
==================
Total nodes: $(jq '.nodes | length' manifest.json)
Total files: $(jq '.nodes[0].rollup.total_files' manifest.json)
Total size: $(jq '.nodes[0].rollup.size.total' manifest.json) bytes
Largest file: $(jq '.nodes[0].rollup.size.max' manifest.json) bytes
Generated: $(jq -r '.generated_at' manifest.json)
EOF
```

---

## Tips and Tricks

### Pretty Print YAML/JSON

**YAML:**
```bash
yq '.' manifest.yaml
```

**JSON:**
```bash
jq '.' manifest.json
```

---

### Validate Manifest Structure

**YAML:**
```bash
# Check if manifest is valid YAML
yq eval 'true' manifest.yaml && echo "✅ Valid YAML" || echo "❌ Invalid YAML"
```

**JSON:**
```bash
# Check if manifest is valid JSON
jq empty manifest.json && echo "✅ Valid JSON" || echo "❌ Invalid JSON"
```

---

### Use with Less for Large Outputs

```bash
yq '.nodes[] | select(.rollup.total_files > 10)' manifest.yaml | less
jq '.nodes[] | select(.rollup.total_files > 10)' manifest.json | less
```

---

## Common Patterns for LLM Workflows

### Generate File List for Specific Task

```bash
# Find all configuration files
yq '.nodes[] | select(.path | test("config|settings|\\.env|\\.yaml|\\.json")) | .path' manifest.yaml

# Find all test files
yq '.nodes[] | select(.path | test("_test\\.go|test_.*\\.py|\\.test\\.")) | .path' manifest.yaml
```

---

### Identify Hot Spots

```bash
# Find directories with most churn (most files)
yq '.nodes[] | select(.rollup.total_files > 20) | {path, files: .rollup.total_files}' manifest.yaml
```

---

### Security/Audit Queries

```bash
# Find potential secret files
yq '.nodes[] | select(.path | test("secret|credential|password|key|token")) | .path' manifest.yaml

# Find large binary files
yq '.nodes[] | select(.size_bytes > 5000000 and (.path | test("\\.bin|\\.exe|\\.so|\\.dylib"))) | {path, size: .size_bytes}' manifest.yaml
```

---

## Further Reading

- [yq documentation](https://mikefarah.gitbook.io/yq/)
- [jq manual](https://jqlang.github.io/jq/manual/)
- [Manifestor README](../README.md)
