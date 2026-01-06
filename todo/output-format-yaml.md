# Proposal: Add YAML Output Format (and Make it Default)

**Status:** Proposed
**Priority:** P1 (High signal, low additional cost)
**Estimated Effort:** ~30 lines of Go code + config update
**Created:** 2026-01-05

---

## Summary

Add YAML as an output format option for manifestor, and make it the **new default** format going forward. Retain JSON support for users who need jq compatibility.

---

## Rationale

### 1. LLM Token Efficiency (Primary Goal)
Manifestor is explicitly designed as an "LLM-friendly manifest" system. YAML delivers:
- **20-30% fewer tokens** vs JSON for typical manifest structures
- **Cleaner syntax:** No `{`, `}`, `"` overhead
- **Better signal-to-noise ratio:** Indentation is semantically meaningful
- **Natural readability:** LLMs are heavily trained on YAML (configs, k8s, CI/CD)

**Example comparison for a single node:**

```json
{
  "path": ".gitignore",
  "is_dir": false,
  "inode": 67038951,
  "mtime_unix": 1767123958,
  "size_bytes": 28
}
```
**~40 tokens**

```yaml
- path: .gitignore
  is_dir: false
  inode: 67038951
  mtime_unix: 1767123958
  size_bytes: 28
```
**~30 tokens** (25% savings)

For a manifest with 100+ nodes, this compounds to **thousands of tokens saved**.

---

### 2. Alignment with Design Philosophy

From the README:
> "This project provides a manifest-first indexing layer for local filesystems, designed specifically for use with LLMs"

YAML directly supports this mission:
- Faster LLM processing (fewer tokens to parse)
- Lower API costs for users
- Better cognitive load for human reviewers
- Still fully machine-readable

---

### 3. Tooling Trade-offs (Acknowledged)

**Concern:** `jq` is more universal than `yq`

**Reality:**
- `jq` is indeed ubiquitous and battle-tested
- `yq` has multiple implementations with varying syntax
- DevOps/platform engineers (target audience) increasingly have `yq` installed
- For ad-hoc queries, users can still opt into JSON

**Mitigation:**
- Support **both** JSON and YAML
- Let users choose based on workflow
- Document common `yq` patterns in README

---

## Proposed Implementation

### Phase 1: Add YAML Support

**Config change:**
```yaml
output:
  format: "yaml"  # or "json"
  file: "manifest.yaml"
```

**Code change (internal/output/writer.go):**
```go
func Write(manifest *Manifest, format string, filepath string) error {
    switch format {
    case "json":
        return writeJSON(manifest, filepath)
    case "yaml":
        return writeYAML(manifest, filepath)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}

func writeYAML(manifest *Manifest, filepath string) error {
    data, err := yaml.Marshal(manifest)
    if err != nil {
        return fmt.Errorf("marshal yaml: %w", err)
    }
    return os.WriteFile(filepath, data, 0644)
}
```

**Dependency:**
```go
import "gopkg.in/yaml.v3"
```

**Default config.yaml:**
```yaml
output:
  format: "yaml"  # NEW DEFAULT
  file: "manifest.yaml"
```

---

### Phase 2: Update Documentation

**README.md addition:**
```markdown
## Output Formats

Manifestor supports multiple output formats to balance LLM optimization with tooling compatibility:

- **YAML** (default): 20-30% fewer tokens for LLM consumption, human-readable
- **JSON**: Maximum tool compatibility (jq, Python, etc.)

### Querying Manifests

**YAML (using yq):**
```bash
# Find directories larger than 10MB
yq '.nodes[] | select(.rollup.size.total > 10000000) | .path' manifest.yaml

# Count files by extension
yq '.nodes[] | select(.is_dir == false) | .path' manifest.yaml | \
  awk -F. '{print $NF}' | sort | uniq -c
```

**JSON (using jq):**
```bash
# Find directories larger than 10MB
jq '.nodes[] | select(.rollup.size.total > 10000000) | .path' manifest.json

# Get top 10 largest files
jq '.nodes[] | select(.size_bytes) | {path, size: .size_bytes}' manifest.json | \
  jq -s 'sort_by(.size) | reverse | .[0:10]'
```

### Switching Formats

```yaml
# config.yaml
output:
  format: "json"  # Use JSON if you prefer jq
  file: "manifest.json"
```
```

---

### Phase 3: Migration Guide for Existing Users

**Breaking change notice (if making YAML default):**
```markdown
## v0.3 Migration Guide

**Breaking change:** Default output format is now YAML instead of JSON.

**If you rely on JSON:**
Update your `config.yaml`:
```yaml
output:
  format: "json"
  file: "manifest.json"
```

**Why the change:**
YAML reduces token usage by 20-30% for LLM consumption, aligning with manifestor's primary design goal. JSON remains fully supported for users who need jq compatibility.
```

---

## Success Metrics

After implementation, verify:

1. **Token count reduction:** Convert an existing manifest.json to YAML and compare token counts
2. **User workflow validation:** Ensure common queries work with both `yq` and `jq`
3. **Performance:** Verify YAML marshaling doesn't add meaningful latency
4. **Backward compatibility:** Existing JSON consumers can opt-in without code changes

---

## Open Questions

1. **Should we support both formats simultaneously?** (e.g., write both manifest.json and manifest.yaml)
   - **Recommendation:** No - let users choose one via config to avoid confusion

2. **Should we add a `--format` CLI flag to override config?**
   - **Recommendation:** Yes, for quick testing/experimentation

3. **Should we add format auto-detection based on file extension?**
   - **Recommendation:** Not initially - explicit is better than implicit

---

## Alternative Considered: TOON Format

**What is TOON?**
- Token-Oriented Object Notation
- Achieves ~40% token reduction vs JSON
- Uses CSV-style tabular layout for uniform arrays

**Why not TOON?**
- Adds external dependency
- Less universal tooling (no toon-query equivalent of jq/yq)
- Learning curve for users
- YAML gets 70% of the benefit with 10% of the risk

**Revisit if:**
- Users explicitly request more aggressive token optimization
- TOON ecosystem matures with better tooling

---

## Next Steps

1. Implement `writeYAML` function in `internal/output/`
2. Add `gopkg.in/yaml.v3` dependency
3. Update default config.yaml to use YAML
4. Add format examples to README
5. Test with existing manifests and measure token savings
6. Document migration path for existing JSON users

---

## Related Work

- `/Users/daniel.nitsch/ais/projects/manifestor/internal/output/json.go` - Existing JSON writer
- `/Users/daniel.nitsch/ais/projects/manifestor/config.yaml` - Current config format
- `/Users/daniel.nitsch/ais/projects/manifestor/README.md` - Documentation updates needed

---

## References

- YAML spec: https://yaml.org/spec/1.2.2/
- Go YAML library: https://github.com/go-yaml/yaml
- yq tool: https://github.com/mikefarah/yq
- TOON format (alternative): https://github.com/toon-format/toon
