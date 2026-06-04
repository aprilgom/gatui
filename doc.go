// Package gatui is a Go port of Ratatui's terminal UI model.
//
// Gatui mirrors Ratatui's separation between backend-independent rendering
// primitives and terminal backends. Core packages such as layout, style, text,
// buffer, symbols, widgets, and terminal should remain usable without an
// interactive terminal.
//
// For Ratatui-to-Gatui porting, start with docs/ratatui-correspondence.md.
// That document preserves Rust paths such as ratatui::widgets::Paragraph and
// maps them to Gatui packages and types. If a Ratatui path is not listed there,
// inspect local Gatui code and tests before assuming the API has been ported.
//
// Gatui intentionally keeps input polling outside core. Applications should
// read keyboard, mouse, and resize events with tcell or another input library,
// then use Gatui terminal and backend packages for drawing.
package gatui
