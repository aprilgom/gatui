package layout

import "fmt"

type Constraint struct {
	kind        constraintKind
	value       int
	denominator int
}

type constraintKind int

const (
	constraintLength constraintKind = iota
	constraintMin
	constraintPercentage
	constraintRatio
	constraintMax
	constraintFill
)

func Length(value int) Constraint {
	return Constraint{kind: constraintLength, value: value}
}

func Min(value int) Constraint {
	return Constraint{kind: constraintMin, value: value}
}

func Max(value int) Constraint {
	return Constraint{kind: constraintMax, value: value}
}

func Percentage(percent int) Constraint {
	return Constraint{kind: constraintPercentage, value: percent}
}

func Ratio(numerator, denominator int) Constraint {
	return Constraint{kind: constraintRatio, value: numerator, denominator: denominator}
}

func Fill(weight int) Constraint {
	return Constraint{kind: constraintFill, value: weight}
}

func DefaultConstraint() Constraint {
	return Percentage(100)
}

func (c Constraint) Apply(length int) int {
	length = maxInt(0, length)
	switch c.kind {
	case constraintPercentage:
		return minInt(int(float64(c.value)/100*float64(length)), length)
	case constraintRatio:
		denominator := maxInt(c.denominator, 1)
		return minInt(int(float64(c.value)/float64(denominator)*float64(length)), length)
	case constraintLength, constraintFill, constraintMax:
		return minInt(length, maxInt(0, c.value))
	case constraintMin:
		return maxInt(length, maxInt(0, c.value))
	default:
		return length
	}
}

func (c Constraint) String() string {
	switch c.kind {
	case constraintLength:
		return fmt.Sprintf("Length(%d)", c.value)
	case constraintMin:
		return fmt.Sprintf("Min(%d)", c.value)
	case constraintPercentage:
		return fmt.Sprintf("Percentage(%d)", c.value)
	case constraintRatio:
		return fmt.Sprintf("Ratio(%d, %d)", c.value, c.denominator)
	case constraintMax:
		return fmt.Sprintf("Max(%d)", c.value)
	case constraintFill:
		return fmt.Sprintf("Fill(%d)", c.value)
	default:
		return fmt.Sprintf("Constraint(%d)", c.kind)
	}
}

func FromLengths(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Length(value)
	}
	return constraints
}

func FromFills(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Fill(value)
	}
	return constraints
}

func FromMaxes(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Max(value)
	}
	return constraints
}

func FromMins(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Min(value)
	}
	return constraints
}

func FromPercentages(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Percentage(value)
	}
	return constraints
}

func FromRatios(values ...[2]int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Ratio(value[0], value[1])
	}
	return constraints
}

func (c Constraint) IsLength() bool {
	return c.kind == constraintLength
}

func (c Constraint) IsMin() bool {
	return c.kind == constraintMin
}

func (c Constraint) IsPercentage() bool {
	return c.kind == constraintPercentage
}

func (c Constraint) IsRatio() bool {
	return c.kind == constraintRatio
}

func (c Constraint) IsMax() bool {
	return c.kind == constraintMax
}

func (c Constraint) IsFill() bool {
	return c.kind == constraintFill
}

func (c Constraint) Value() int {
	return c.value
}

func (c Constraint) Denominator() int {
	return c.denominator
}
