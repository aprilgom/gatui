package layout

func axisLength(area Rect, direction Direction) int {
	if direction == Vertical {
		return area.Height
	}
	return area.Width
}

func spacerRect(area Rect, direction Direction, start int, length int) Rect {
	if direction == Vertical {
		return Rect{X: area.X, Y: area.Y + start, Width: area.Width, Height: length}
	}
	return Rect{X: area.X + start, Y: area.Y, Width: length, Height: area.Height}
}

func emptySpacer(area Rect, direction Direction, start int) Rect {
	return spacerRect(area, direction, start, 0)
}

func calculateLengths(areaLength int, constraints []Constraint, stretchFixedSurplus bool) []int {
	areaLength = maxInt(0, areaLength)
	lengths := make([]int, len(constraints))
	totalFixed := 0
	minIndexes := make([]int, 0)
	totalPositiveFillWeight := 0
	fillCount := 0

	for i, constraint := range constraints {
		if constraint.kind == constraintFill {
			fillCount++
			if constraint.value > 0 {
				totalPositiveFillWeight += constraint.value
			}
			continue
		}

		length := constraintLengthValue(areaLength, constraint)
		lengths[i] = length
		totalFixed += length
		if constraint.kind == constraintMin {
			minIndexes = append(minIndexes, i)
		}
	}
	if allProportionalConstraints(constraints) {
		totalFixed += distributeProportionalRemainders(areaLength, lengths, constraints)
	}

	if fillCount > 0 {
		if totalFixed > areaLength {
			shrinkLengths(lengths, constraints, totalFixed-areaLength, false)
			shrinkLengths(lengths, constraints, sumInts(lengths)-areaLength, true)
			return lengths
		}

		distributeFillLengths(lengths, constraints, areaLength-totalFixed, totalPositiveFillWeight)
		return lengths
	}

	total := sumInts(lengths)
	switch {
	case total < areaLength:
		surplus := areaLength - total
		if len(minIndexes) > 0 {
			distributeSurplus(lengths, minIndexes, surplus)
		} else if stretchFixedSurplus && len(lengths) > 0 {
			lengths[len(lengths)-1] += surplus
		}
	case total > areaLength:
		if allProportionalConstraints(constraints) {
			shrinkLargestLengths(lengths, total-areaLength)
		} else {
			shrinkLengths(lengths, constraints, total-areaLength, false)
			shrinkLengths(lengths, constraints, sumInts(lengths)-areaLength, true)
		}
	}

	return lengths
}

func flexOffsets(areaLength int, lengths []int, flex Flex, spacing int) []int {
	offsets := make([]int, len(lengths))
	if len(lengths) == 0 {
		return offsets
	}

	total := spacedLength(lengths, spacing)
	surplus := maxInt(0, areaLength-total)
	switch flex {
	case FlexEnd:
		setPackedOffsets(offsets, lengths, surplus, spacing)
	case FlexCenter:
		setPackedOffsets(offsets, lengths, roundedDiv(surplus, 2), spacing)
	case FlexSpaceBetween:
		if len(lengths) == 1 {
			offsets[0] = 0
			return offsets
		}
		if spacing < 0 {
			setPackedOffsets(offsets, lengths, 0, spacing)
			return offsets
		}
		surplus = maxInt(0, areaLength-sumInts(lengths))
		for i := range lengths {
			offsets[i] = sumInts(lengths[:i]) + roundedDiv(i*surplus, len(lengths)-1)
		}
	case FlexSpaceAround:
		surplus = maxInt(0, areaLength-sumInts(lengths))
		denominator := len(lengths) * 2
		for i := range lengths {
			offsets[i] = sumInts(lengths[:i]) + roundedDiv((2*i+1)*surplus, denominator)
		}
	case FlexSpaceEvenly:
		surplus = maxInt(0, areaLength-sumInts(lengths))
		denominator := len(lengths) + 1
		for i := range lengths {
			offsets[i] = sumInts(lengths[:i]) + roundedDiv((i+1)*surplus, denominator)
		}
	default:
		setPackedOffsets(offsets, lengths, 0, spacing)
	}

	return offsets
}

func setPackedOffsets(offsets []int, lengths []int, leading int, spacing int) {
	cursor := leading
	for i, length := range lengths {
		offsets[i] = cursor
		cursor += length + spacing
	}
}

func spacedLength(lengths []int, spacing int) int {
	total := sumInts(lengths)
	if len(lengths) > 1 {
		total += spacing * (len(lengths) - 1)
	}
	return total
}

func distributeSurplus(lengths []int, indexes []int, surplus int) {
	if surplus <= 0 || len(indexes) == 0 {
		return
	}

	base := surplus / len(indexes)
	remainder := surplus % len(indexes)
	for _, index := range indexes {
		lengths[index] += base
		if remainder > 0 {
			lengths[index]++
			remainder--
		}
	}
}

func roundedDiv(numerator, denominator int) int {
	if denominator <= 0 {
		return 0
	}
	return (numerator + denominator/2) / denominator
}

func constraintLengthValue(areaLength int, constraint Constraint) int {
	switch constraint.kind {
	case constraintLength:
		return clampInt(constraint.value, 0, areaLength)
	case constraintMin:
		return clampInt(constraint.value, 0, areaLength)
	case constraintMax:
		return clampInt(constraint.value, 0, areaLength)
	case constraintPercentage:
		percent := clampInt(constraint.value, 0, 100)
		return areaLength * percent / 100
	case constraintRatio:
		if constraint.denominator <= 0 {
			return areaLength
		}
		return clampInt(areaLength*constraint.value/constraint.denominator, 0, areaLength)
	default:
		return 0
	}
}

func centeredLength(areaLength int, constraint Constraint) int {
	lengths := calculateLengths(areaLength, []Constraint{constraint}, false)
	if len(lengths) == 0 {
		return 0
	}
	return minInt(areaLength, lengths[0])
}

func allProportionalConstraints(constraints []Constraint) bool {
	if len(constraints) == 0 {
		return false
	}
	for _, constraint := range constraints {
		if constraint.kind != constraintPercentage && constraint.kind != constraintRatio {
			return false
		}
	}
	return true
}

func distributeProportionalRemainders(areaLength int, lengths []int, constraints []Constraint) int {
	type remainder struct {
		index int
		value float64
	}

	target := 0.0
	remainders := make([]remainder, 0, len(constraints))
	for i, constraint := range constraints {
		exact := 0.0
		switch constraint.kind {
		case constraintPercentage:
			exact = float64(areaLength*clampInt(constraint.value, 0, 100)) / 100
		case constraintRatio:
			if constraint.denominator <= 0 {
				exact = float64(areaLength)
			} else {
				exact = float64(areaLength*constraint.value) / float64(constraint.denominator)
			}
		}
		target += exact
		remainders = append(remainders, remainder{index: i, value: exact - float64(lengths[i])})
	}

	extra := roundedDiv(int(target*1000), 1000) - sumInts(lengths)
	added := 0
	for ; extra > 0; extra-- {
		best := -1
		for i, remainder := range remainders {
			if remainder.value <= 0 {
				continue
			}
			if best == -1 || remainder.value > remainders[best].value {
				best = i
			}
		}
		if best == -1 {
			return added
		}
		lengths[remainders[best].index]++
		remainders[best].value = 0
		added++
	}
	return added
}

func shrinkLargestLengths(lengths []int, shortage int) {
	for shortage > 0 {
		best := -1
		for i, length := range lengths {
			if length == 0 {
				continue
			}
			if best == -1 || length > lengths[best] {
				best = i
			}
		}
		if best == -1 {
			return
		}
		lengths[best]--
		shortage--
	}
}

func distributeFillLengths(lengths []int, constraints []Constraint, remaining int, totalPositiveWeight int) {
	if remaining <= 0 {
		return
	}

	if totalPositiveWeight <= 0 {
		fillCount := 0
		for _, constraint := range constraints {
			if constraint.kind == constraintFill {
				fillCount++
			}
		}
		if fillCount == 0 {
			return
		}

		base := remaining / fillCount
		remainder := remaining % fillCount
		for i, constraint := range constraints {
			if constraint.kind != constraintFill {
				continue
			}
			lengths[i] = base
			if remainder > 0 {
				lengths[i]++
				remainder--
			}
		}
		return
	}

	distributed := 0
	type fillRemainder struct {
		index     int
		remainder int
	}
	remainders := make([]fillRemainder, 0)

	for i, constraint := range constraints {
		if constraint.kind != constraintFill || constraint.value <= 0 {
			continue
		}

		scaled := remaining * constraint.value
		length := scaled / totalPositiveWeight
		lengths[i] = length
		distributed += length
		remainders = append(remainders, fillRemainder{index: i, remainder: scaled % totalPositiveWeight})
	}

	for leftover := remaining - distributed; leftover > 0; leftover-- {
		best := 0
		for i := 1; i < len(remainders); i++ {
			if remainders[i].remainder > remainders[best].remainder {
				best = i
			}
		}
		lengths[remainders[best].index]++
		remainders[best].remainder = 0
	}
}

func shrinkLengths(lengths []int, constraints []Constraint, shortage int, includeMin bool) {
	for i := len(lengths) - 1; i >= 0 && shortage > 0; i-- {
		if constraints[i].kind == constraintMin && !includeMin {
			continue
		}
		reduction := minInt(lengths[i], shortage)
		lengths[i] -= reduction
		shortage -= reduction
	}
}

func sumInts(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampInt(value, low, high int) int {
	if high < low {
		return low
	}
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
