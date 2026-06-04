package symbols

type BarSet struct {
	Empty         string
	OneEighth     string
	OneQuarter    string
	ThreeEighths  string
	Half          string
	FiveEighths   string
	ThreeQuarters string
	SevenEighths  string
	Full          string
}

var NineLevelBarSet = BarSet{
	Empty:         " ",
	OneEighth:     "▁",
	OneQuarter:    "▂",
	ThreeEighths:  "▃",
	Half:          "▄",
	FiveEighths:   "▅",
	ThreeQuarters: "▆",
	SevenEighths:  "▇",
	Full:          "█",
}

var ThreeLevelBarSet = BarSet{
	Empty:         " ",
	OneEighth:     "▄",
	OneQuarter:    "▄",
	ThreeEighths:  "▄",
	Half:          "▄",
	FiveEighths:   "█",
	ThreeQuarters: "█",
	SevenEighths:  "█",
	Full:          "█",
}
