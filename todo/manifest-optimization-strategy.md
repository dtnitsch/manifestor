# Manifest Optimization Strategy

**Status:** Proposed
**Priority:** P1 (High value, medium effort)
**Created:** 2026-01-05
**Estimated Impact:** **~34% token reduction** (~1,020 tokens saved on current manifest)

---

## Executive Summary

Three complementary optimizations can reduce manifest token usage by **~34%** while maintaining:
- ✅ Human readability
- ✅ jq/yq queryability
- ✅ Lossless representation
- ✅ No decompression step

**Current manifest.json:** ~2,989 tokens (474 lines, 10.5KB)
**Optimized manifest.yaml:** ~1,969 tokens (estimated ~320 lines, ~7KB)
**Savings:** ~1,020 tokens (34.1%)

For a typical 3,000-line manifest (~6,000 tokens), this translates to **~2,047 tokens saved**.

---

## Token Savings Breakdown

### Phase 1: Quick Wins (v0.3) - **31% savings**

#### ✅ Optimization 1: Omit Null/Empty/Default Values
**Token savings:** ~243 tokens (8.1%)

**Before:**
```json
{
  "path": ".gitignore",
  "is_dir": false,
  "inode": 67038951,
  "mtime_unix": 1767123958,
  "size_bytes": 28,
  "direct_subdir_count": null,
  "rollup": null
}
```
**~80 tokens**

**After:**
```json
{
  "path": ".gitignore",
  "inode": 67038951,
  "mtime_unix": 1767123958,
  "size_bytes": 28
}
```
**~55 tokens (31% reduction per node)**

**Implementation:**
```go
type Node struct {
    Path              string   `json:"path" yaml:"path"`
    IsDir             bool     `json:"is_dir,omitempty" yaml:"is_dir,omitempty"`
    Inode             *uint64  `json:"inode,omitempty" yaml:"inode,omitempty"`
    MtimeUnix         *int64   `json:"mtime_unix,omitempty" yaml:"mtime_unix,omitempty"`
    SizeBytes         *int64   `json:"size_bytes,omitempty" yaml:"size_bytes,omitempty"`
    DirectSubdirCount *int     `json:"direct_subdir_count,omitempty" yaml:"direct_subdir_count,omitempty"`
    Rollup            *Rollup  `json:"rollup,omitempty" yaml:"rollup,omitempty"`
}
```

**Rules:**
- Omit `is_dir: false` (files are default)
- Omit `rollup: null`
- Omit `direct_subdir_count: 0` or `null`
- Omit empty arrays and objects
- Only include fields with meaningful values

**Effort:** 1-2 hours (struct tag changes + pointer conversions)
**Breaking change:** No (missing fields treated as zero/null/false by consumers)

---

#### ✅ Optimization 2: Switch to YAML
**Additional token savings:** ~686 tokens (23.0%)
**Cumulative savings:** ~930 tokens (31.1%)

**After omitempty (JSON):**
```json
{
  "path": ".gitignore",
  "inode": 67038951,
  "mtime_unix": 1767123958,
  "size_bytes": 28
}
```
**~55 tokens**

**After YAML:**
```yaml
- path: .gitignore
  inode: 67038951
  mtime_unix: 1767123958
  size_bytes: 28
```
**~40 tokens (27% reduction per node)**

**Why YAML saves tokens:**
- No `{`, `}`, `,` punctuation
- No quotes around keys
- Indentation instead of braces
- More compact array syntax

**Effort:** 2-3 hours (see `todo/output-format-yaml.md`)
**Breaking change:** Yes (but JSON remains supported)

---

### Phase 2: Structural Optimization (v0.4) - **Additional 3% savings**

#### ✅ Optimization 3: Delta-Encoded Timestamps
**Additional token savings:** ~90 tokens (3.0%)
**Cumulative savings:** ~1,020 tokens (34.1%)

**Current (absolute timestamps):**
```yaml
- path: file1.go
  mtime_unix: 1767379653  # 10 digits = ~3 tokens
- path: file2.go
  mtime_unix: 1767379658  # 10 digits = ~3 tokens
```

**Optimized (delta from generation time):**
```yaml
generated_at: 1767379653
nodes:
  - path: file1.go
    mt_delta: 0           # 1 digit = ~1 token
  - path: file2.go
    mt_delta: 5           # 1 digit = ~1 token
```

**Calculation:**
- 45 timestamps in current manifest
- Average savings: ~2 tokens per timestamp
- Total: ~90 tokens saved

**Implementation:**
```go
func (s *Scanner) calculateTimeDelta(mtime int64) int64 {
    return mtime - s.baseTime
}

// During manifest generation:
manifest.GeneratedAt = time.Now().Unix()
for _, node := range nodes {
    if node.MtimeUnix != nil {
        delta := *node.MtimeUnix - manifest.GeneratedAt
        node.MtimeDelta = &delta
        node.MtimeUnix = nil  // Remove absolute time
    }
}
```

**jq/yq reconstruction helper (add to README):**
```bash
# Reconstruct absolute timestamps
BASE=$(yq '.generated_at' manifest.yaml)
yq ".nodes[] | {path, mtime: (.mt_delta + $BASE)}" manifest.yaml
```

**Effort:** 4-5 hours (field changes + query documentation)
**Breaking change:** Yes (field renamed, requires calculation)
**Decision gate:** Measure timestamp clustering in real-world manifests first

---

## Combined Impact: Real Numbers

### Current Manifest (474 lines)

| Format | Characters | Lines | Tokens | Savings |
|--------|-----------|-------|--------|---------|
| **Current JSON** | 10,461 | 474 | ~2,989 | — |
| JSON + omitempty | 9,610 | ~430 | ~2,746 | 243 (8.1%) |
| YAML + omitempty | 7,208 | ~320 | ~2,059 | 930 (31.1%) |
| **YAML + omitempty + delta** | ~6,890 | ~310 | **~1,969** | **1,020 (34.1%)** |

### Extrapolated: Large Manifest (3,000 lines)

| Format | Estimated Tokens | Savings |
|--------|-----------------|---------|
| **Current JSON** | ~6,000 | — |
| JSON + omitempty | ~5,514 | 486 (8.1%) |
| YAML + omitempty | ~4,134 | 1,866 (31.1%) |
| **YAML + omitempty + delta** | **~3,953** | **2,047 (34.1%)** |

---

## Cost Impact Analysis

### Token Costs (Claude 3.5 Sonnet pricing)

**Input tokens:** $3.00 per million tokens

| Scenario | Manifest Size | Cost per 1,000 reads | Annual cost (daily reads) |
|----------|---------------|---------------------|---------------------------|
| Current JSON | 2,989 tokens | $8.97 | $3,274 |
| Optimized YAML | 1,969 tokens | $5.91 | $2,157 |
| **Savings** | **1,020 tokens** | **$3.06** | **$1,117/year** |

For large repos with 6,000-token manifests read daily by LLMs:
- **Current cost:** $6,570/year
- **Optimized cost:** $4,330/year
- **Savings:** $2,240/year

---

## Implementation Plan

### Phase 1: v0.3 (Immediate - 3-5 hours)

**Week 1:**
1. Add `omitempty` tags to all structs
2. Convert fields to pointers where needed
3. Add YAML output support (see `output-format-yaml.md`)
4. Update default config to use YAML
5. Test with existing manifests

**Deliverables:**
- ~930 tokens saved (31.1%)
- Backward-compatible (JSON still supported)
- No query syntax changes

**Validation:**
```bash
# Generate both formats
./manifestor --format json -o manifest.json
./manifestor --format yaml -o manifest.yaml

# Compare sizes
echo "JSON tokens: $(cat manifest.json | wc -c | awk '{print int($1/3.5)}')"
echo "YAML tokens: $(cat manifest.yaml | wc -c | awk '{print int($1/3.5)}')"

# Verify query compatibility
jq '.nodes[] | select(.is_dir == true)' manifest.json > /tmp/json-output
yq '.nodes[] | select(.is_dir == true)' manifest.yaml > /tmp/yaml-output
diff /tmp/json-output /tmp/yaml-output && echo "✅ Queries match"
```

---

### Phase 2: v0.4 (After validation - 4-6 hours)

**Prerequisites:**
- Measure timestamp clustering in 5+ real repositories
- Verify savings justify complexity
- Document query helpers

**Week 2-3:**
1. Add `mt_delta` field to Node struct
2. Calculate deltas during scan
3. Update output serialization
4. Add query examples to README
5. Update manifest version to 0.4

**Deliverables:**
- Additional ~90 tokens saved (3% more)
- Total savings: 34.1%
- Query syntax changes documented

---

## Rejected Approaches

### ❌ Short Key Names
```json
{"p": "file.go", "mt": 1767379653, "sz": 1024}
```
**Why rejected:** Destroys human readability, breaks intuitive jq queries

### ❌ Hierarchical Paths
```yaml
internal:
  manifest:
    - capabilities.go
    - node.go
```
**Why rejected:** Changes data model, breaks flat-list queries, higher complexity

### ❌ TOON Format
**Why rejected:** Gets 40% savings vs 34% for YAML, but adds:
- External dependency
- New tooling learning curve
- Less universal support
- Not worth the additional 6% for the complexity

### ❌ Compact Extension Maps
```json
"extensions": [[".go", 42], [".md", 3]]
```
**Why rejected:** Only 2-5% savings, breaks intuitive queries, not worth complexity

---

## Testing Strategy

### Token Counting Script
```bash
#!/bin/bash
# test/token-comparison.sh

echo "=== Generating test manifests ==="
./manifestor --format json -o /tmp/test.json
./manifestor --format yaml -o /tmp/test.yaml

echo ""
echo "=== Size Comparison ==="
json_bytes=$(wc -c < /tmp/test.json)
yaml_bytes=$(wc -c < /tmp/test.yaml)
json_tokens=$((json_bytes * 10 / 35))  # chars / 3.5
yaml_tokens=$((yaml_bytes * 10 / 35))

echo "JSON: $json_bytes bytes (~$json_tokens tokens)"
echo "YAML: $yaml_bytes bytes (~$yaml_tokens tokens)"
echo "Savings: $((json_tokens - yaml_tokens)) tokens ($((100 * (json_tokens - yaml_tokens) / json_tokens))%)"

echo ""
echo "=== Query Compatibility ==="
jq '.nodes | length' /tmp/test.json > /tmp/json-count
yq '.nodes | length' /tmp/test.yaml > /tmp/yaml-count
if diff /tmp/json-count /tmp/yaml-count > /dev/null; then
    echo "✅ Node counts match"
else
    echo "❌ Node counts differ - data loss detected!"
    exit 1
fi

echo "✅ All tests passed"
```

### Regression Tests
```go
// internal/manifest/optimization_test.go

func TestOmitEmptyFields(t *testing.T) {
    node := &Node{
        Path:              "test.go",
        IsDir:             false,  // Should be omitted
        DirectSubdirCount: nil,    // Should be omitted
        SizeBytes:         ptr(int64(100)),
    }

    data, _ := json.Marshal(node)

    // Verify is_dir is not in output
    assert.NotContains(t, string(data), "is_dir")

    // Verify size_bytes IS in output
    assert.Contains(t, string(data), "size_bytes")
}

func TestDeltaTimestamps(t *testing.T) {
    baseTime := int64(1767379653)
    mtime := int64(1767379658)

    delta := calculateTimeDelta(mtime, baseTime)
    assert.Equal(t, int64(5), delta)

    // Verify reconstruction
    reconstructed := baseTime + delta
    assert.Equal(t, mtime, reconstructed)
}
```

---

## Success Metrics

**After Phase 1 (v0.3):**
- ✅ Token reduction: ≥30%
- ✅ jq/yq queries: 100% compatible (no syntax changes)
- ✅ Human readability: Maintained or improved
- ✅ Generation time: <5% increase
- ✅ File size: 30-35% smaller

**After Phase 2 (v0.4):**
- ✅ Token reduction: ≥33%
- ✅ jq/yq queries: 95%+ compatible (timestamp queries need helper)
- ✅ Documentation: Query examples in README
- ✅ Manifest version: Bumped to 0.4

---

## Migration Guide

### For Existing JSON Users

**No action required if you want to keep JSON:**
```yaml
# config.yaml
output:
  format: "json"
  file: "manifest.json"
```

**To adopt YAML (recommended):**
```yaml
# config.yaml
output:
  format: "yaml"
  file: "manifest.yaml"
```

Install `yq` (if not already installed):
```bash
# macOS
brew install yq

# Linux
wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq
chmod +x /usr/bin/yq
```

**Common query translations:**
```bash
# Find large directories
# JSON:
jq '.nodes[] | select(.rollup.size.total > 10000000)' manifest.json
# YAML:
yq '.nodes[] | select(.rollup.size.total > 10000000)' manifest.yaml

# Count files by extension
# JSON:
jq '.nodes[].rollup.extensions | keys[]' manifest.json | sort | uniq -c
# YAML:
yq '.nodes[].rollup.extensions | keys[]' manifest.yaml | sort | uniq -c
```

---

## Open Questions

1. **Should we support dual output (both JSON and YAML)?**
   - **Recommendation:** No - increases complexity, users can choose one

2. **Should delta timestamps be optional via config?**
   - **Recommendation:** Yes - add `optimizations.delta_timestamps: true/false`

3. **How to version these changes?**
   - **Recommendation:** Bump manifest version for each phase:
     - v0.3: omitempty + YAML support
     - v0.4: delta timestamps

4. **Should we warn users about breaking changes?**
   - **Recommendation:** Yes - add migration guide to README and `--check-version` flag

---

## Future Optimizations (v0.5+)

### Folder-Based Splitting
Your original idea of splitting into multiple files:
```
manifest/
  summary.yaml       # Always loaded (~200 tokens)
  files.yaml         # Load on demand (~1,500 tokens)
  rollups.yaml       # Load on demand (~270 tokens)
```

**This is orthogonal to token optimization** - it's about **selective loading** rather than compression.

**Recommendation:** Pursue this separately as a "lazy loading" feature, not a compression strategy.

---

## Related Work

- `internal/manifest/node.go` - Node struct definitions
- `internal/manifest/rollup.go` - Rollup struct definitions
- `internal/output/json.go` - JSON marshaling
- `todo/output-format-yaml.md` - YAML output proposal
- `config.yaml` - Default configuration

---

## References

- Claude tokenizer: ~3.5 chars per token (measured)
- YAML spec: https://yaml.org/spec/1.2.2/
- Go struct tags: https://pkg.go.dev/encoding/json#Marshal
- yq documentation: https://mikefarah.gitbook.io/yq/

---

## Conclusion

**Phase 1 (omitempty + YAML) delivers 31% token savings** with minimal complexity and zero loss of queryability. This should be implemented immediately.

**Phase 2 (delta timestamps) adds another 3%** but requires more careful consideration of the query complexity trade-off.

**Total potential: 34% token reduction (~1,020 tokens on current manifest, ~2,047 on large repos)**

For an LLM-focused tool like manifestor, this directly translates to:
- Faster processing
- Lower API costs
- Better user experience
- Alignment with stated mission

**Recommended action:** Implement Phase 1 in v0.3, evaluate Phase 2 based on real-world usage data.
