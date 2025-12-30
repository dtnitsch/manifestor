## Negative Learning (Things We Learned the Hard Way)

This project exists largely because a number of *obvious-sounding approaches* consistently failed in practice.

If you’re working with LLMs and local codebases, these mistakes are easy to make — and expensive to recover from.

---

### 1. “Just read the whole repo”

This is the most common failure mode.

Problems:

* Token limits are hit immediately
* Generated and vendor code dominate context
* Important files are drowned out by volume
* Results are brittle and non-reproducible

Even small repos become unusable once history, build output, or dependencies are included.

**Lesson:**
Reading is the *last* step, not the first.

---

### 2. Grep-first workflows don’t scale to reasoning

Tools like `grep`, `ripgrep`, and `find` are excellent for humans — but misleading for LLMs.

They:

* return unranked matches
* flatten structure
* remove semantic context
* bias toward frequently repeated strings

An LLM needs *relative importance*, not raw matches.

**Lesson:**
Search is not understanding.

---

### 3. File size and recency are poor proxies for importance

We repeatedly assumed:

* bigger files matter more
* newer files are more relevant

In practice:

* the most important files are often small
* core logic is stable and rarely touched
* configuration and documentation punch far above their size

**Lesson:**
Importance is structural and semantic, not temporal.

---

### 4. Jumping straight to embeddings hides problems

Embedding entire repos:

* is expensive
* is slow
* is hard to invalidate
* locks you into a single retrieval strategy

Worse, it hides poor ingestion decisions behind vector similarity.

**Lesson:**
If you can’t explain *why* something was indexed, you indexed too early.

---

### 5. “LLM as filesystem explorer” is the wrong mental model

Giving an LLM tool access to:

* `ls`
* `cat`
* recursive reads

…encourages:

* unbounded exploration
* excessive reads
* context bloat
* non-deterministic behavior

This works in demos and fails in real systems.

**Lesson:**
LLMs should *decide what to read*, not wander until they hit limits.

---

### 6. Mixing indexing and reasoning makes everything worse

Early prototypes combined:

* filesystem walking
* content reading
* summarization
* reasoning

into a single pass.

This made:

* caching impossible
* debugging painful
* behavior unpredictable

**Lesson:**
Indexing must be deterministic. Reasoning can be probabilistic.

---

### 7. Re-scanning everything kills iteration speed

Without stable manifests:

* every prompt becomes a cold start
* nothing is reusable
* costs grow linearly with repo size

This makes LLM-assisted workflows feel slow and fragile.

**Lesson:**
If it can be cached, it should be cached.

---

### 8. LLMs perform better with *structure than content*

Counterintuitively, LLM output quality improved when we gave:

* fewer raw files
* more structural summaries
* clearer constraints

Even when the LLM had *less* total data.

**Lesson:**
Clarity beats completeness.

---

## Why this section exists

These mistakes aren’t theoretical — they’re the default path.

This project exists to encode these lessons into infrastructure, so you don’t have to rediscover them under token limits, time pressure, or broken demos.
