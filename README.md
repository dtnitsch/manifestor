# Manifest-First Context for LLMs

**Stop asking LLMs to read your filesystem.
Start telling them what matters.**

## What this is

This project provides a **manifest-first indexing layer** for local filesystems, designed specifically for use with **LLMs and agentic tools**.

Instead of dumping folders, grepping blindly, or re-reading massive repos on every prompt, this system produces a **deterministic, incremental, and LLM-friendly manifest** that describes your codebase *without reading its contents*.

The manifest becomes the stable interface between:

* your filesystem
* LLMs (ChatGPT, Claude Code, Gemini, Ollama, etc.)
* downstream reasoning, summarization, and analysis tools

---

## Why you should care

### 1. “Read this folder” does not scale

LLMs struggle when given:

* thousands of files
* large generated directories
* mixed-quality content
* unclear intent

Most tools jump straight from **existence → ingestion**.

This project inserts a missing step: **understanding before reading**.

---

### 2. You get answers faster — with fewer tokens

The manifest lets an LLM:

* see repo structure
* identify hotspots
* rank files by importance
* decide *what not to read*

Result:

* fewer tokens
* less noise
* more accurate reasoning
* dramatically better follow-up questions

---

### 3. Works with *any* LLM — local or hosted

This is not tied to a specific model.

The manifest is:

* plain text (YAML / JSON)
* deterministic
* versionable
* cacheable

That makes it ideal for:

* ChatGPT (no filesystem access)
* Claude Code / Gemini CLI
* local Ollama agents
* CI or offline analysis

The LLM doesn’t need access to your files — it reasons over the manifest.

---

### 4. Enables “surgical” LLM workflows

Instead of:

> “Here’s my repo, can you look at it?”

You get:

> “Please review `src/auth/tokens.go` and `src/config/secrets.yaml`. Ignore everything else.”

This flips the workflow:

* LLMs **plan**
* tools **execute**
* context stays bounded and intentional

---

### 5. Incremental by design

The manifest uses:

* structural hashes
* file metadata
* directory-level change detection

Which means:

* no full re-scans
* no re-reading unchanged trees
* no wasted compute

It behaves more like a **filesystem index** than a crawler.

---

## What the manifest contains (high level)

Without reading file contents, the manifest captures:

* directory structure
* file counts and sizes
* language and type signals
* “likely generated” vs “hand-written”
* complexity and size hotspots
* top files by size / line count / word count
* stable hashes for change detection

Think:
`ls`, `stat`, and human intuition — formalized.

---

## What this unlocks

Once you have a manifest, you can:

* ask an LLM *what to read first*
* generate repo summaries safely
* guide refactors
* focus security reviews
* power RAG pipelines without blind ingestion
* cache and reuse understanding across sessions

The manifest becomes the **control plane for context**.

---

## Design philosophy

* **Manifest before meaning**
* **Rank before read**
* **Summarize before ingest**
* **Deterministic over clever**
* **LLMs decide, tools execute**

This is infrastructure — not a prompt trick.

---

## Who this is for

* Engineers using LLMs on large repos
* Platform / infra teams
* Security reviewers
* Tooling authors
* Anyone tired of “just paste the code”

If you’ve ever thought:

> “The model would be great if it just knew where to look”

This is for you.

---

## Status

Early, but intentionally designed to be:

* small
* composable
* LLM-agnostic
* CLI-friendly

Proof comes next.
