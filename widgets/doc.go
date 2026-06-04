// Package widgets corresponds to ratatui::widgets.
//
// It contains renderable UI components such as Block, Paragraph, Clear, Gauge,
// LineGauge, BarChart, Sparkline, List, Table, Tabs, Scrollbar, Canvas,
// Calendar, Fill, and Shadow. Widget interfaces in this package are the Go
// equivalents of Ratatui widget traits where they have been ported.
//
// Use this package when porting Ratatui code that imports
// ratatui::widgets::{Widget, StatefulWidget, Block, Paragraph} or concrete
// widget types. Translate Rust builder chains to the constructors, fields, and
// methods exported by the Gatui widget type.
package widgets
