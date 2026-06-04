// Package terminal corresponds to Ratatui's Terminal, Frame, Viewport, and
// backend orchestration concepts.
//
// Terminal owns drawing over a backend, while Frame represents the rendering
// context for a single draw pass. Backend is the Go interface implemented by
// concrete terminal backends and test backends.
//
// Gatui does not mirror ratatui::init, ratatui::restore, or ratatui::run as a
// complete application event loop. Applications should create a backend, draw
// through Terminal, and handle input polling with tcell or another input
// library outside core.
package terminal
