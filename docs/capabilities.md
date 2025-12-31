# Manifest Capabilities

This document defines the **capabilities model** for `manifestor`, with a focus on **rollup semantics**, **validation guarantees**, and **forward compatibility**.

Capabilities describe **what a generated manifest promises to contain**. Validation logic uses these declarations to assert invariants and detect partial, inconsistent, or corrupted rollups.

---

## Design Goals

* **Explicit contracts**: A manifest must declare what it claims to support.
* **Cheap validation**: Capability checks must be fast and stateless.
* **Forward compatible**: Older consumers can safely ignore unknown capabilities.
* **Single-pass derivation**: All rollup capabilities must be computable during a single filesystem scan.

---

## Capability Structure

Capabilities are declared at the manifest level:

```json
{
  "manifest": {
    "capabilities": {
      "rollup": { ... }
    }
  }
}
```

Currently, only **rollup capabilities** are defined.

---

## RollupCapabilities

```go
type RollupCapabilities struct {
    SizeStats        bool `json:"size_stats"`
    SizePercentiles  bool `json:"size_percentiles"`
    SizeBuckets      bool `json:"size_buckets"`
    ActivitySpan     bool `json:"activity_span"`

    ExtensionCounts  bool `json:"extension_counts"`
    DirCounts        bool `json:"dir_counts"`

    DepthMetrics     bool `json:"depth_metrics"`
}
```

Each flag indicates that **all directory rollups** in the manifest satisfy the corresponding invariant set.

---

## Capability Definitions & Invariants

### `size_stats`

Rollups include aggregate file size statistics.

**Required fields:**

* `rollup.size.total`
* `rollup.size.min`
* `rollup.size.max`
* `rollup.size.mean`
* `rollup.size.median`

**Invariants:**

* `total > 0` if `total_files > 0`
* `min <= median <= max`

---

### `size_percentiles`

Rollups include percentile statistics derived from file sizes.

**Required fields:**

* `rollup.size.percentiles.p50`
* `rollup.size.percentiles.p90`
* `rollup.size.percentiles.p99`

**Invariants:**

* `median == p50`
* `min <= p50 <= p90 <= p99 <= max`

---

### `size_buckets`

Rollups include coarse-grained size distribution buckets.

**Required fields:**

* `rollup.size.buckets`

**Notes:**

* Bucket definitions are fixed per schema version.
* Intended for fast binary/blob detection.

---

### `activity_span`

Rollups include modification-time span metrics.

**Required fields:**

* `rollup.last_modified`
* (future) `oldest_mtime`, `newest_mtime`, `span_seconds`

---

### `extension_counts`

Rollups include counts of file extensions.

**Required fields:**

* `rollup.extensions`

**Invariants:**

* Sum of extension counts == `total_files`

---

### `dir_counts`

Rollups include directory count aggregation.

**Required fields:**

* `rollup.total_descendant_dirs`

**Invariants:**

* `total_descendant_dirs >= direct_subdir_count`

---

### `depth_metrics` (planned)

Rollups include directory tree depth statistics.

**Status:** Not yet implemented.

---

## Validation Model

Capabilities are validated during manifest construction via:

* `validateCapabilitiesInline()` (current)
* A future **table-driven invariant engine**

Only capabilities explicitly set to `true` are validated.

---

## Compatibility Rules

* Consumers **must not assume** a capability unless declared.
* Producers **must not declare** a capability unless all invariants hold.
* Unknown capabilities must be ignored by consumers.

---

## Roadmap

* Introduce a capability invariant table
* Promote select capabilities from experimental → stable
* Expose capability summaries to downstream LLM ingestion pipelines

---

## Philosophy

Capabilities are not features.
They are **promises**.

A manifest that lies about its capabilities is worse than one that omits them entirely.

