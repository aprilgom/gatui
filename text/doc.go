// Package text corresponds to ratatui::text.
//
// It defines styled text primitives such as Span, Line, Text, Masked, and
// StyledGrapheme where those concepts are exported by Gatui. Use this package
// when porting Ratatui code that imports ratatui::text::{Span, Line, Text}.
//
// Text values remain backend-independent and can be rendered into Gatui buffers
// or consumed by widgets.
package text
