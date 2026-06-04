package symbols

type ScrollbarSet struct {
	Track string
	Thumb string
	Begin string
	End   string
}

var (
	HorizontalScrollbarSet = ScrollbarSet{
		Track: "═",
		Thumb: "█",
		Begin: "◄",
		End:   "►",
	}
	VerticalScrollbarSet = ScrollbarSet{
		Track: "║",
		Thumb: "█",
		Begin: "▲",
		End:   "▼",
	}
)
