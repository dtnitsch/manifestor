# Manifestor – Design Overview

## Purpose

Manifestor scans a local filesystem and produces a structured **manifest** describing
folders and files without reading file contents.

The goal is to create a **stable, bounded, machine-readable representation** of a filesystem
that can be used by humans, tools, or LLMs for reasoning, summarization, and follow-on analysis.

This project explicitly avoids content ingestion in v0.x.

---

## Core Concepts

### Manifest

A manifest is a hierarchical summary of a filesystem rooted at a single directory.

At minimum, it captures:
- directory structure
- file names
- metadata (timestamps, counts, inode info where available)

Future versions may enrich this, but the manifest is always:
- deterministic
- reproducible
- content-agnostic

---

### Scanner

The scanner is responsible for:
- walking the filesystem
- applying allow/block rules
- collecting metadata
- emitting a manifest

The scanner:
- uses Go-native recursive traversal
- executes work in parallel using a bounded worker pool
- does **not** follow symlinks (by design)

---

### Filters

Filters determine which paths are included or excluded during scanning.

Rules are:
- explicit
- ordered
- deterministic

Block rules are evaluated first.
Allow rules can override block rules.

This avoids glob explosion, hidden magic, and unbounded traversal.

---

## High-Level Flow

1. Load configuration
2. Initialize scanner with options and filters
3. Recursively walk filesystem
4. Apply filtering rules per path
5. Collect metadata
6. Emit manifest
7. Write output

---

## Non-Goals (v0.x)

Manifestor intentionally does **not**:
- read file contents
- tokenize files
- parse programming languages
- follow symlinks
- watch for filesystem changes
- provide a server or API

These are possible future extensions, but are explicitly out of scope for now.

---

## Design Principles

- Prefer clarity over cleverness
- Favor deterministic output
- Avoid implicit behavior
- Keep debugging simple
- Treat configuration as a user-facing artifact

---

## Why This Exists

Most tools treat filesystem discovery as an implementation detail.

Manifestor treats it as a **first-class problem** — one worth modeling carefully,
especially when the output is consumed by automated systems or LLMs.

This project exists to make filesystem structure *understandable*, not just enumerable.

