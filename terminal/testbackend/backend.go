package testbackend

import (
	"reflect"
	"strings"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/terminal"
)

type Backend struct {
	size             layout.Size
	draws            [][]buffer.CellDiff
	flushCount       int
	clearCount       int
	clearRegions     []terminal.ClearType
	hideCursorCount  int
	showCursorCount  int
	cursorPositions  []layout.Position
	cursorPosition   layout.Position
	cursorVisible    bool
	appendLines      []int
	scrollback       []string
	scrollRegionUp   [][3]int
	scrollRegionDown [][3]int
	cells            *buffer.Buffer
}

type NoScrollBackend struct {
	backend *Backend
}

func New(width, height int) *Backend {
	area := layout.NewRect(0, 0, width, height)
	return &Backend{size: layout.Size{Width: width, Height: height}, cursorVisible: true, cells: buffer.Empty(area)}
}

func NewNoScroll(width, height int) *NoScrollBackend {
	return &NoScrollBackend{backend: New(width, height)}
}

func WithLines(lines []string) *Backend {
	cells := buffer.WithLines(lines)
	return &Backend{
		size:          layout.Size{Width: cells.Area.Width, Height: cells.Area.Height},
		cursorVisible: true,
		cells:         cells,
	}
}

func WithLinesNoScroll(lines []string) *NoScrollBackend {
	return &NoScrollBackend{backend: WithLines(lines)}
}

func (b *NoScrollBackend) Size() (layout.Size, error) {
	return b.backend.Size()
}

func (b *NoScrollBackend) WindowSize() (terminal.WindowSize, error) {
	return b.backend.WindowSize()
}

func (b *NoScrollBackend) SetSize(width, height int) {
	b.backend.SetSize(width, height)
}

func (b *NoScrollBackend) Draw(diffs []buffer.CellDiff) error {
	return b.backend.Draw(diffs)
}

func (b *NoScrollBackend) Flush() error {
	return b.backend.Flush()
}

func (b *NoScrollBackend) Clear() error {
	return b.backend.Clear()
}

func (b *NoScrollBackend) ClearRegion(clearType terminal.ClearType) error {
	return b.backend.ClearRegion(clearType)
}

func (b *NoScrollBackend) HideCursor() error {
	return b.backend.HideCursor()
}

func (b *NoScrollBackend) ShowCursor() error {
	return b.backend.ShowCursor()
}

func (b *NoScrollBackend) SetCursorPosition(pos layout.Position) error {
	return b.backend.SetCursorPosition(pos)
}

func (b *NoScrollBackend) GetCursorPosition() (layout.Position, error) {
	return b.backend.GetCursorPosition()
}

func (b *NoScrollBackend) AppendLines(count int) error {
	return b.backend.AppendLines(count)
}

func (b *NoScrollBackend) Lines() []string {
	return b.backend.Lines()
}

func (b *NoScrollBackend) Buffer() *buffer.Buffer {
	return b.backend.Buffer()
}

func (b *NoScrollBackend) String() string {
	return b.backend.String()
}

func (b *NoScrollBackend) Scrollback() *buffer.Buffer {
	return b.backend.Scrollback()
}

func (b *NoScrollBackend) AssertBufferLines(t testing.TB, lines []string) {
	b.backend.AssertBufferLines(t, lines)
}

func (b *NoScrollBackend) AssertScrollbackLines(t testing.TB, lines []string) {
	b.backend.AssertScrollbackLines(t, lines)
}

func (b *NoScrollBackend) AssertScrollbackEmpty(t testing.TB) {
	b.backend.AssertScrollbackEmpty(t)
}

func (b *NoScrollBackend) AssertCursorPosition(t testing.TB, pos layout.Position) {
	b.backend.AssertCursorPosition(t, pos)
}

func (b *NoScrollBackend) AppendLinesCalls() []int {
	return b.backend.AppendLinesCalls()
}

func (b *Backend) Size() (layout.Size, error) {
	return b.size, nil
}

func (b *Backend) WindowSize() (terminal.WindowSize, error) {
	return terminal.WindowSize{
		ColumnsRows: b.size,
		Pixels:      layout.Size{Width: 640, Height: 480},
	}, nil
}

func (b *Backend) SetSize(width, height int) {
	b.size = layout.Size{Width: width, Height: height}
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, width, height))
		return
	}
	b.cells.Resize(layout.NewRect(0, 0, width, height))
}

func (b *Backend) Draw(diffs []buffer.CellDiff) error {
	copied := make([]buffer.CellDiff, len(diffs))
	copy(copied, diffs)
	b.draws = append(b.draws, copied)
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	for _, diff := range diffs {
		b.cells.SetCell(diff.X, diff.Y, diff.Cell)
	}
	return nil
}

func (b *Backend) Flush() error {
	b.flushCount++
	return nil
}

func (b *Backend) Clear() error {
	b.clearCount++
	return b.ClearRegion(terminal.ClearAll)
}

func (b *Backend) ClearRegion(clearType terminal.ClearType) error {
	b.clearRegions = append(b.clearRegions, clearType)
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	switch clearType {
	case terminal.ClearAll:
		b.cells.Reset()
	case terminal.ClearAfterCursor:
		for y := b.cursorPosition.Y; y < b.size.Height; y++ {
			startX := 0
			if y == b.cursorPosition.Y {
				startX = b.cursorPosition.X
			}
			for x := startX; x < b.size.Width; x++ {
				b.cells.SetCell(x, y, buffer.NewCell(" "))
			}
		}
	case terminal.ClearBeforeCursor:
		for y := 0; y <= b.cursorPosition.Y && y < b.size.Height; y++ {
			endX := b.size.Width - 1
			if y == b.cursorPosition.Y {
				endX = b.cursorPosition.X
			}
			for x := 0; x <= endX && x < b.size.Width; x++ {
				b.cells.SetCell(x, y, buffer.NewCell(" "))
			}
		}
	case terminal.ClearCurrentLine:
		for x := 0; x < b.size.Width; x++ {
			b.cells.SetCell(x, b.cursorPosition.Y, buffer.NewCell(" "))
		}
	case terminal.ClearUntilNewLine:
		for x := b.cursorPosition.X; x < b.size.Width; x++ {
			b.cells.SetCell(x, b.cursorPosition.Y, buffer.NewCell(" "))
		}
	}
	return nil
}

func (b *Backend) HideCursor() error {
	b.hideCursorCount++
	b.cursorVisible = false
	return nil
}

func (b *Backend) ShowCursor() error {
	b.showCursorCount++
	b.cursorVisible = true
	return nil
}

func (b *Backend) SetCursorPosition(pos layout.Position) error {
	b.cursorPositions = append(b.cursorPositions, pos)
	b.cursorPosition = pos
	return nil
}

func (b *Backend) GetCursorPosition() (layout.Position, error) {
	return b.cursorPosition, nil
}

func (b *Backend) AppendLines(count int) error {
	b.appendLines = append(b.appendLines, count)
	if count <= 0 {
		return nil
	}
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	newX := b.cursorPosition.X + 1
	if newX >= b.size.Width {
		newX = b.size.Width - 1
	}
	maxY := b.size.Height - 1
	linesAfterCursor := maxY - b.cursorPosition.Y
	if linesAfterCursor < 0 {
		linesAfterCursor = 0
	}

	if count > linesAfterCursor {
		scroll := count - linesAfterCursor
		visibleScroll := scroll
		if visibleScroll > b.size.Height {
			visibleScroll = b.size.Height
		}
		for y := 0; y < visibleScroll; y++ {
			b.scrollback = append(b.scrollback, b.lineAt(y))
		}
		for y := 0; y < b.size.Height-visibleScroll; y++ {
			for x := 0; x < b.size.Width; x++ {
				cell, _ := b.cells.CellAt(x, y+visibleScroll)
				b.cells.SetCell(x, y, cell)
			}
		}
		for y := b.size.Height - visibleScroll; y < b.size.Height; y++ {
			b.clearLine(y)
		}
		for y := visibleScroll; y < scroll; y++ {
			b.scrollback = append(b.scrollback, b.blankLine())
		}
	}

	newY := b.cursorPosition.Y + count
	if newY > maxY {
		newY = maxY
	}
	b.cursorPosition = layout.Position{X: newX, Y: newY}
	return nil
}

func (b *Backend) ScrollRegionUp(startY, endY, count int) error {
	b.scrollRegionUp = append(b.scrollRegionUp, [3]int{startY, endY, count})
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	if count <= 0 || startY >= endY {
		return nil
	}
	if startY < 0 {
		startY = 0
	}
	if endY > b.size.Height {
		endY = b.size.Height
	}
	height := endY - startY
	if height <= 0 {
		return nil
	}
	if count > height {
		count = height
	}
	if startY == 0 {
		for y := 0; y < count; y++ {
			b.scrollback = append(b.scrollback, b.lineAt(y))
		}
	}
	for y := startY; y < endY-count; y++ {
		for x := 0; x < b.size.Width; x++ {
			cell, _ := b.cells.CellAt(x, y+count)
			b.cells.SetCell(x, y, cell)
		}
	}
	for y := endY - count; y < endY; y++ {
		b.clearLine(y)
	}
	return nil
}

func (b *Backend) ScrollRegionDown(startY, endY, count int) error {
	b.scrollRegionDown = append(b.scrollRegionDown, [3]int{startY, endY, count})
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	if count <= 0 || startY >= endY {
		return nil
	}
	if startY < 0 {
		startY = 0
	}
	if endY > b.size.Height {
		endY = b.size.Height
	}
	height := endY - startY
	if height <= 0 {
		return nil
	}
	if count > height {
		count = height
	}
	for y := endY - 1; y >= startY+count; y-- {
		for x := 0; x < b.size.Width; x++ {
			cell, _ := b.cells.CellAt(x, y-count)
			b.cells.SetCell(x, y, cell)
		}
	}
	for y := startY; y < startY+count; y++ {
		b.clearLine(y)
	}
	return nil
}

func (b *Backend) lineAt(y int) string {
	if b.cells == nil || y < 0 || y >= b.size.Height {
		return ""
	}
	line := make([]rune, 0, b.size.Width)
	for x := 0; x < b.size.Width; x++ {
		cell, _ := b.cells.CellAt(x, y)
		symbol := cell.DisplaySymbol()
		for _, r := range symbol {
			line = append(line, r)
		}
	}
	return string(line)
}

func (b *Backend) clearLine(y int) {
	for x := 0; x < b.size.Width; x++ {
		b.cells.SetCell(x, y, buffer.NewCell(" "))
	}
}

func (b *Backend) blankLine() string {
	return strings.Repeat(" ", b.size.Width)
}

func (b *Backend) Draws() [][]buffer.CellDiff {
	draws := make([][]buffer.CellDiff, len(b.draws))
	for i := range b.draws {
		draws[i] = make([]buffer.CellDiff, len(b.draws[i]))
		copy(draws[i], b.draws[i])
	}
	return draws
}

func (b *Backend) FlushCount() int {
	return b.flushCount
}

func (b *Backend) ClearCount() int {
	return b.clearCount
}

func (b *Backend) ClearRegions() []terminal.ClearType {
	return append([]terminal.ClearType(nil), b.clearRegions...)
}

func (b *Backend) Lines() []string {
	if b.cells == nil {
		return nil
	}
	return b.cells.Lines()
}

func (b *Backend) String() string {
	return bufferView(b.Buffer())
}

func (b *Backend) Buffer() *buffer.Buffer {
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	return b.cells
}

func bufferView(buf *buffer.Buffer) string {
	if buf == nil || buf.Area.Height == 0 {
		return ""
	}
	var builder strings.Builder
	for y := buf.Area.Y; y < buf.Area.Y+buf.Area.Height; y++ {
		builder.WriteByte('"')
		for x := buf.Area.X; x < buf.Area.X+buf.Area.Width; x++ {
			cell, _ := buf.CellAt(x, y)
			builder.WriteString(cell.DisplaySymbol())
		}
		builder.WriteString("\"\n")
	}
	return builder.String()
}

func (b *Backend) Scrollback() *buffer.Buffer {
	return buffer.WithLines(b.scrollback)
}

func (b *Backend) AssertBufferLines(t testing.TB, lines []string) {
	t.Helper()
	got := b.Lines()
	if !reflect.DeepEqual(got, lines) {
		t.Fatalf("buffer lines mismatch:\nactual: %#v\nwant:   %#v", got, lines)
	}
}

func (b *Backend) AssertScrollbackLines(t testing.TB, lines []string) {
	t.Helper()
	got := b.ScrollbackLines()
	if !reflect.DeepEqual(got, lines) {
		t.Fatalf("scrollback lines mismatch:\nactual: %#v\nwant:   %#v", got, lines)
	}
}

func (b *Backend) AssertScrollbackEmpty(t testing.TB) {
	t.Helper()
	b.AssertScrollbackLines(t, nil)
}

func (b *Backend) HideCursorCount() int {
	return b.hideCursorCount
}

func (b *Backend) ShowCursorCount() int {
	return b.showCursorCount
}

func (b *Backend) CursorPositions() []layout.Position {
	return append([]layout.Position(nil), b.cursorPositions...)
}

func (b *Backend) CursorVisible() bool {
	return b.cursorVisible
}

func (b *Backend) CursorPosition() layout.Position {
	return b.cursorPosition
}

func (b *Backend) AssertCursorPosition(t testing.TB, pos layout.Position) {
	t.Helper()
	if b.cursorPosition != pos {
		t.Fatalf("cursor position mismatch: actual: %+v, want: %+v", b.cursorPosition, pos)
	}
}

func (b *Backend) AppendLinesCalls() []int {
	return append([]int(nil), b.appendLines...)
}

func (b *Backend) ScrollbackLines() []string {
	return append([]string(nil), b.scrollback...)
}

func (b *Backend) ScrollRegionUpCalls() [][3]int {
	return append([][3]int(nil), b.scrollRegionUp...)
}

func (b *Backend) ScrollRegionDownCalls() [][3]int {
	return append([][3]int(nil), b.scrollRegionDown...)
}
