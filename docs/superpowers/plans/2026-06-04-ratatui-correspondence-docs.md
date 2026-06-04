# Ratatui Correspondence Docs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add search-friendly Ratatui-to-Gatui documentation for LLMs, agents, and Go users.

**Architecture:** Create one canonical Markdown correspondence map and focused package-level Go doc files. Keep the change documentation-only and avoid altering behavior or API names.

**Tech Stack:** Go documentation comments, Markdown, `go test`, `go doc`.

---

### Task 1: Canonical Correspondence Map

**Files:**
- Create: `docs/ratatui-correspondence.md`

- [ ] **Step 1: Add the Markdown map**

Create `docs/ratatui-correspondence.md` with:

```markdown
# Ratatui to Gatui Correspondence

This document maps Ratatui Rust API paths to Gatui Go packages and types.
```

- [ ] **Step 2: Include porting rules**

Add agent rules explaining how to search for `ratatui::...` paths, prefer mapped Gatui APIs, and treat missing mappings as unported unless local code proves otherwise.

- [ ] **Step 3: Include package and type mappings**

Add mappings for root concepts, `layout`, `style`, `text`, `buffer`, `symbols`, `widgets`, `terminal`, and `backend/tcell`.

### Task 2: Go Package Documentation

**Files:**
- Create: `doc.go`
- Create: `layout/doc.go`
- Create: `style/doc.go`
- Create: `text/doc.go`
- Create: `buffer/doc.go`
- Create: `symbols/doc.go`
- Create: `widgets/doc.go`
- Create: `terminal/doc.go`
- Create: `backend/tcell/doc.go`

- [ ] **Step 1: Add root package documentation**

Create root `doc.go` with package-level Gatui orientation and Ratatui correspondence guidance.

- [ ] **Step 2: Add package-level correspondence docs**

Create one `doc.go` in each package with the matching `package` declaration and a concise Ratatui path mapping.

### Task 3: Verification

**Files:**
- Read: all created files

- [ ] **Step 1: Run tests**

Run:

```sh
go test ./...
```

Expected: all packages pass.

- [ ] **Step 2: Inspect generated Go docs**

Run:

```sh
go doc .
go doc ./widgets
go doc ./terminal
```

Expected: package summaries include Ratatui correspondence guidance.
