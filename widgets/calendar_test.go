package widgets_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/widgets"
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

func TestMonthly_widthReflectsGridLayout(t *testing.T) {
	monthly := widgets.NewMonthly(utcDate(2015, time.February, 1), nil)
	if got := monthly.Width(); got != 21 {
		t.Fatalf("Width() = %d, want 21", got)
	}

	monthly = monthly.Block(widgets.BorderedBlock().Padding(widgets.NewPadding(2, 3, 1, 2)))
	if got := monthly.Width(); got != 28 {
		t.Fatalf("Width() with block = %d, want 28", got)
	}
}

func TestMonthly_heightCountsWeeksAndHeaders(t *testing.T) {
	monthly := widgets.NewMonthly(utcDate(2015, time.February, 1), nil)
	if got := monthly.Height(); got != 4 {
		t.Fatalf("Height() = %d, want 4", got)
	}

	monthly = monthly.
		ShowMonthHeader(style.NewStyle()).
		ShowWeekdaysHeader(style.NewStyle())
	if got := monthly.Height(); got != 6 {
		t.Fatalf("Height() with headers = %d, want 6", got)
	}
}

func TestMonthly_dimensionsExamples(t *testing.T) {
	monthly := widgets.NewMonthly(utcDate(2015, time.February, 1), nil).
		ShowMonthHeader(style.NewStyle()).
		ShowWeekdaysHeader(style.NewStyle()).
		Block(widgets.BorderedBlock().Padding(widgets.NewPadding(2, 3, 1, 2)))

	if got := monthly.Width(); got != 28 {
		t.Fatalf("Width() = %d, want 28", got)
	}
	if got := monthly.Height(); got != 11 {
		t.Fatalf("Height() = %d, want 11", got)
	}
}

func TestCalendarEventStore_today(t *testing.T) {
	events := widgets.NewCalendarEventStore()
	eventStyle := style.NewStyle().Fg(style.Red)
	before := time.Now()

	events.Today(eventStyle)

	after := time.Now()
	for _, now := range []time.Time{before, after} {
		key := localDateKey(now)
		if got, ok := events[key]; ok {
			if got != eventStyle {
				t.Fatalf("Today() style for %s = %#v, want %#v", key, got, eventStyle)
			}
			return
		}
	}
	t.Fatalf("Today() did not add local today key; checked %s and %s", localDateKey(before), localDateKey(after))
}

func utcDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func localDateKey(date time.Time) string {
	local := date.Local()
	return fmt.Sprintf("%04d-%02d-%02d", local.Year(), local.Month(), local.Day())
}
