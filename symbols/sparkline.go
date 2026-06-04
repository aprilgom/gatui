package symbols

type SparklineBarSet struct {
	Full          string
	SevenEighths  string
	ThreeQuarters string
	FiveEighths   string
	Half          string
	ThreeEighths  string
	OneQuarter    string
	OneEighth     string
	Empty         string
}

func NineLevelSparklineBarSet() SparklineBarSet {
	return SparklineBarSet{
		Full:          "█",
		SevenEighths:  "▇",
		ThreeQuarters: "▆",
		FiveEighths:   "▅",
		Half:          "▄",
		ThreeEighths:  "▃",
		OneQuarter:    "▂",
		OneEighth:     "▁",
		Empty:         " ",
	}
}

func ThreeLevelSparklineBarSet() SparklineBarSet {
	return SparklineBarSet{
		Full:          "█",
		SevenEighths:  "█",
		ThreeQuarters: "▄",
		FiveEighths:   "▄",
		Half:          "▄",
		ThreeEighths:  "▄",
		OneQuarter:    "▄",
		OneEighth:     " ",
		Empty:         " ",
	}
}

func (s SparklineBarSet) SymbolForHeight(height uint64) string {
	switch height {
	case 0:
		return s.Empty
	case 1:
		return s.OneEighth
	case 2:
		return s.OneQuarter
	case 3:
		return s.ThreeEighths
	case 4:
		return s.Half
	case 5:
		return s.FiveEighths
	case 6:
		return s.ThreeQuarters
	case 7:
		return s.SevenEighths
	default:
		return s.Full
	}
}
