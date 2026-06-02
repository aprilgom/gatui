package layout

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

func FromFills(values ...int) []Constraint {
	constraints := make([]Constraint, len(values))
	for i, value := range values {
		constraints[i] = Fill(value)
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
