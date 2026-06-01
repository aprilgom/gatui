package widgets_test

import (
	"testing"
	"time"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/widgets"
)

func TestMonthly_shouldRenderDaysLayout(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 5))

	widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  1  2  3  4  5  6  7",
		"  8  9 10 11 12 13 14",
		" 15 16 17 18 19 20 21",
		" 22 23 24 25 26 27 28",
		" 29 30 31            ",
	})
}

func TestMonthly_shouldRenderSurroundingDays(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 6))

	widgets.NewMonthly(utcDate(2023, time.December, 1), widgets.NewCalendarEventStore()).
		ShowSurrounding(style.NewStyle()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		" 26 27 28 29 30  1  2",
		"  3  4  5  6  7  8  9",
		" 10 11 12 13 14 15 16",
		" 17 18 19 20 21 22 23",
		" 24 25 26 27 28 29 30",
		" 31  1  2  3  4  5  6",
	})
}

func TestMonthly_shouldRenderMonthHeader(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 6))

	widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
		ShowMonthHeader(style.NewStyle()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    January 2023     ",
		"  1  2  3  4  5  6  7",
		"  8  9 10 11 12 13 14",
		" 15 16 17 18 19 20 21",
		" 22 23 24 25 26 27 28",
		" 29 30 31            ",
	})
}

func TestMonthly_shouldRenderWeekdaysHeader(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 6))

	widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
		ShowWeekdaysHeader(style.NewStyle()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		" Su Mo Tu We Th Fr Sa",
		"  1  2  3  4  5  6  7",
		"  8  9 10 11 12 13 14",
		" 15 16 17 18 19 20 21",
		" 22 23 24 25 26 27 28",
		" 29 30 31            ",
	})
}

func TestMonthly_shouldRenderCombinedHeadersAndSurrounding(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 7))

	widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
		ShowWeekdaysHeader(style.NewStyle()).
		ShowMonthHeader(style.NewStyle()).
		ShowSurrounding(style.NewStyle()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    January 2023     ",
		" Su Mo Tu We Th Fr Sa",
		"  1  2  3  4  5  6  7",
		"  8  9 10 11 12 13 14",
		" 15 16 17 18 19 20 21",
		" 22 23 24 25 26 27 28",
		" 29 30 31  1  2  3  4",
	})
}

func TestMonthly_shouldPatchEventStyleOnDateCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 21, 5))
	events := widgets.NewCalendarEventStore()
	events.Add(utcDate(2023, time.January, 10), style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold))

	widgets.NewMonthly(utcDate(2023, time.January, 22), events).
		DefaultStyle(style.NewStyle().Bg(style.Blue)).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  1  2  3  4  5  6  7",
		"  8  9 10 11 12 13 14",
		" 15 16 17 18 19 20 21",
		" 22 23 24 25 26 27 28",
		" 29 30 31            ",
	})
	for x := 6; x <= 8; x++ {
		assertCellStyle(t, buf, x, 1, style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold))
	}
}

func TestMonthly_shouldTruncateSmallAreaWithoutPanic(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))

	assertNotPanics(t, func() {
		widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
			ShowMonthHeader(style.NewStyle()).
			Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{
		"Janua",
		"  1  ",
	})
}

func TestMonthly_shouldRenderInsideBlock(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 23, 7))

	widgets.NewMonthly(utcDate(2023, time.January, 1), widgets.NewCalendarEventStore()).
		Block(widgets.BorderedBlock()).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌─────────────────────┐",
		"│  1  2  3  4  5  6  7│",
		"│  8  9 10 11 12 13 14│",
		"│ 15 16 17 18 19 20 21│",
		"│ 22 23 24 25 26 27 28│",
		"│ 29 30 31            │",
		"└─────────────────────┘",
	})
}

func utcDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
