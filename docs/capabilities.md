# Rollup Capabilities Specification

**Project:** scanner  
**Version:** 0.2  
**Status:** Draft (P0 semantics locking)  
**Last updated:** 2025-12-30

This document defines **rollup capabilities** declared in `manifest.capabilities.rollup`
and the **invariants that MUST hold** if a capability is claimed.

Capabilities exist to:
- make rollup semantics explicit
- prevent silent partial data
- enable safe downstream reasoning
- support schema evolution without breakage

If a capability is declared and its invariants are violated, the manifest is **invalid**.

---

## Capability Model

Capabilities are **opt-in guarantees**, not feature flags.

Rules:
- Declaring a capability asserts *presence and correctness* of specific fields
- Absence of a capability means **no guarantees**
- Validators MUST NOT infer capabilities from data presence
- Capabilities apply only to directory nodes with a rollup

---

## P0 Capabilities (v0.2)

### 1. `size_stats`

**Description:**  
Directory-level aggregate file size statistics.

**If declared, the following MUST hold for every directory rollup:**
- `rollup.size.total` is present
- `rollup.size.total > 0` if `rollup.total_files > 0`
- `rollup.size.min` is present
- `rollup.size.max` is present
- `rollup.size.mean` is present
- `rollup.size.median` is present
- `rollup.size.min <= rollup.size.median <= rollup.size.max`

**Notes:**
- Units are bytes
- Statistics are computed over files only (not directories)

---

### 2. `size_buckets`

**Description:**  
File count distribution across coarse size ranges.

**If declared, the following MUST hold:**
- `rollup.size.buckets` is present
- All bucket keys defined by the schema are present
- Sum of all bucket counts == `rollup.total_files`

**Defined buckets (v0.2):**
- `<1KB`
- `1KB-1MB`
- `1MB-10MB`
- `>10MB`

---

### 3. `activity_span`

**Description:**  
Modification time span of files within a directory subtree.

**If declared, the following MUST hold:**
- `rollup.last_modified` is present
- Value represents the **maximum mtime** of all descendant files
- `rollup.last_modified >= node.mtime_unix`

**Notes:**
- Derived from stat() only
- Does not imply continuous activity, only bounds

---

### 4. `capability_integrity` *(meta-capability)*

**Description:**  
Indicates that the manifest enforces capability-driven validation.

**If declared, the following MUST hold:**
- All declared rollup capabilities are validated
- Validation failures are fatal
- Validation errors reference:
  - node path
  - violated capability
  - invariant description

**Notes:**
- This capability applies to the manifest as a whole
- It is implicitly required for production use

---

### 5. `rollup_completeness`

**Description:**  
Indicates whether rollups are complete across directory nodes.

**If declared, the following MUST hold:**
- Every directory node either:
  - has a rollup, or
  - is explicitly marked as excluded
- Manifest includes:
  - `directories_with_rollups`
  - `directories_missing_rollups`

**Notes:**
- Enables consumers to reason about partial scans

---

## Non-Goals (v0.2)

- File content inspection
- MIME detection
- Heuristic scoring
- Cross-manifest comparisons

---

## Schema Evolution Rules (Preview)

- Capabilities are additive
- Removing a capability is a breaking change
- New invariants MUST be guarded by new capabilities
- Validators MUST ignore unknown capabilities

---

## Status

- Capability definitions: **draft**
- Invariants: **locking**
- Validator: **in-progress**

