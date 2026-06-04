// Package layout corresponds to ratatui::layout.
//
// It contains backend-independent geometry and area-splitting primitives such
// as Rect, Margin, Offset, Constraint, Direction, and Layout. Use this package
// when porting Ratatui code that imports ratatui::layout::{Rect, Constraint,
// Layout} or related layout types.
//
// Gatui follows Go naming and API conventions, so Rust builder chains should be
// translated to the constructors, fields, and methods exported by this package.
package layout
