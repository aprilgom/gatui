// Package tcell provides a Gatui terminal backend implemented with tcell.
//
// This package is the closest Gatui counterpart to Ratatui applications using
// crossterm-backed terminal drawing, but it is not a direct port of
// ratatui_crossterm. Gatui applications should use tcell or another input
// library for keyboard, mouse, and resize events, then use this backend for
// drawing buffers to the terminal.
package tcell
