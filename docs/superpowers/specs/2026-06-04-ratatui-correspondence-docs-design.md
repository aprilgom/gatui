# Ratatui Correspondence Docs Design

## Goal

Make Gatui's documentation useful to LLMs, agents, and Go users who already know Ratatui's Rust API. A reader should be able to search for a `ratatui::...` path and find the closest Gatui package, type, or concept.

## Scope

This first pass adds documentation only. It does not change runtime behavior, public API names, or port missing Ratatui features.

## Architecture

Use one canonical correspondence document plus package-level Go documentation:

- `docs/ratatui-correspondence.md` is the search-friendly source of truth for Ratatui-to-Gatui mappings.
- Package `doc.go` files expose the same orientation through Go documentation and pkg.go.dev.
- Root package documentation explains the overall porting model and directs agents to the correspondence map.

## Content Rules

- Preserve Ratatui Rust paths such as `ratatui::widgets::Paragraph` literally so search and LLM retrieval work.
- Prefer "corresponds to" wording for implemented concepts.
- Use "not yet ported" or "applications own this" for missing or intentionally external behavior.
- Explain Go idioms where Ratatui uses Rust builders, traits, or backends.

## Verification

Documentation must compile with `go test ./...`. Go documentation should be inspectable with `go doc` for root and package docs.
