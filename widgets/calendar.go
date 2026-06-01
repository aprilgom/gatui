package widgets

import (
	"fmt"
	"time"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type CalendarEventStore map[string]style.Style

type Monthly struct {
	displayDate         time.Time
	events              CalendarEventStore
	defaultStyle        style.Style
	surroundingStyle    style.Style
	monthHeaderStyle    style.Style
	weekdaysHeaderStyle style.Style
	showSurrounding     bool
	showMonthHeader     bool
	showWeekdaysHeader  bool
	block               *Block
}

func NewCalendarEventStore() CalendarEventStore {
	return CalendarEventStore{}
}

func (s CalendarEventStore) Add(date time.Time, eventStyle style.Style) {
	s[dateKey(date)] = eventStyle
}

func NewMonthly(displayDate time.Time, events CalendarEventStore) Monthly {
	if events == nil {
		events = NewCalendarEventStore()
	}
	return Monthly{
		displayDate:  monthStart(displayDate),
		events:       events,
		defaultStyle: style.NewStyle(),
	}
}

func (m Monthly) ShowSurrounding(surroundingStyle style.Style) Monthly {
	m.showSurrounding = true
	m.surroundingStyle = surroundingStyle
	return m
}

func (m Monthly) ShowWeekdaysHeader(headerStyle style.Style) Monthly {
	m.showWeekdaysHeader = true
	m.weekdaysHeaderStyle = headerStyle
	return m
}

func (m Monthly) ShowMonthHeader(headerStyle style.Style) Monthly {
	m.showMonthHeader = true
	m.monthHeaderStyle = headerStyle
	return m
}

func (m Monthly) DefaultStyle(defaultStyle style.Style) Monthly {
	m.defaultStyle = defaultStyle
	return m
}

func (m Monthly) Block(block Block) Monthly {
	m.block = &block
	return m
}

func (m Monthly) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	calendarArea := area
	if m.block != nil {
		m.block.Render(area, buf)
		calendarArea = m.block.Inner(area)
	}
	if calendarArea.Width == 0 || calendarArea.Height == 0 {
		return
	}
	buf.SetStyle(calendarArea, m.defaultStyle)

	y := calendarArea.Y
	if m.showMonthHeader && y < calendarArea.Y+calendarArea.Height {
		m.renderCentered(calendarArea.X, y, calendarArea.Width, m.displayDate.Format("January 2006"), m.monthHeaderStyle, buf)
		y++
	}
	if m.showWeekdaysHeader && y < calendarArea.Y+calendarArea.Height {
		m.writeString(calendarArea.X, y, calendarArea.X+calendarArea.Width, " Su Mo Tu We Th Fr Sa", m.weekdaysHeaderStyle, buf)
		y++
	}

	m.renderDays(calendarArea.X, y, calendarArea.X+calendarArea.Width, calendarArea.Y+calendarArea.Height, buf)
}

func (m Monthly) renderDays(left, top, right, bottom int, buf *buffer.Buffer) {
	if top >= bottom || left >= right {
		return
	}
	month := m.displayDate.Month()
	nextMonth := monthStart(m.displayDate.AddDate(0, 1, 0))
	date := firstCalendarDate(m.displayDate)
	y := top
	for {
		if y >= bottom {
			return
		}
		for weekday := 0; weekday < 7; weekday++ {
			x := left + weekday*3
			if x >= right {
				continue
			}
			inMonth := date.Month() == month
			if inMonth || m.showSurrounding {
				cellStyle := m.defaultStyle
				if !inMonth {
					cellStyle = cellStyle.Patch(m.surroundingStyle)
				}
				if eventStyle, ok := m.events[dateKey(date)]; ok {
					cellStyle = cellStyle.Patch(eventStyle)
				}
				m.writeString(x, y, right, fmt.Sprintf("%3d", date.Day()), cellStyle, buf)
			}
			date = date.AddDate(0, 0, 1)
		}
		y++
		if !date.Before(nextMonth) && date.Weekday() == time.Sunday {
			return
		}
	}
}

func (m Monthly) renderCentered(x, y, width int, value string, cellStyle style.Style, buf *buffer.Buffer) {
	runes := []rune(value)
	if len(runes) > width {
		runes = runes[:width]
	}
	offset := (width - len(runes)) / 2
	for i, r := range runes {
		buf.SetCell(x+offset+i, y, buffer.Cell{Symbol: string(r), Style: m.defaultStyle.Patch(cellStyle)})
	}
}

func (m Monthly) writeString(x, y, right int, value string, cellStyle style.Style, buf *buffer.Buffer) {
	patchedStyle := m.defaultStyle.Patch(cellStyle)
	for _, r := range value {
		if x >= right {
			return
		}
		buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: patchedStyle})
		x++
	}
}

func firstCalendarDate(date time.Time) time.Time {
	first := monthStart(date)
	return first.AddDate(0, 0, -int(first.Weekday()))
}

func monthStart(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func dateKey(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day())
}
