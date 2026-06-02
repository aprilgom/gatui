package layout

import (
	"fmt"
	"math"
	"sync"

	"github.com/aprilgom/casow"
)

type solverStrength struct {
	value float64
}

var (
	strengthSpacerSize     = solverStrength{value: casow.Required.Div(10).Value()}
	strengthMinSize        = solverStrength{value: casow.Strong.Mul(100).Value()}
	strengthMaxSize        = solverStrength{value: casow.Strong.Mul(100).Value()}
	strengthLengthSize     = solverStrength{value: casow.Strong.Mul(10).Value()}
	strengthPercentageSize = solverStrength{value: casow.Strong.Value()}
	strengthRatioSize      = solverStrength{value: casow.Strong.Div(10).Value()}
	strengthMinSizeEq      = solverStrength{value: casow.Medium.Mul(10).Value()}
	strengthMaxSizeEq      = solverStrength{value: casow.Medium.Mul(10).Value()}
	strengthFillGrow       = solverStrength{value: casow.Medium.Value()}
	strengthGrow           = solverStrength{value: casow.Medium.Div(10).Value()}
	strengthSpaceGrow      = solverStrength{value: casow.Weak.Mul(10).Value()}
	strengthAllSegmentGrow = solverStrength{value: casow.Weak.Value()}
)

func (s solverStrength) casow() casow.Strength {
	return casow.NewStrength(s.value)
}

func (s solverStrength) isValid() bool {
	strength := s.casow()
	return s.value >= 0 && strength.Value() == s.value
}

type solvedLayout struct {
	segments []Rect
	spacers  []Rect
}

type layoutCacheKey struct {
	area        Rect
	direction   Direction
	flex        Flex
	spacing     int
	margin      Margin
	constraints string
}

var layoutCache = struct {
	sync.Mutex
	size  int
	order []layoutCacheKey
	items map[layoutCacheKey]solvedLayout
}{}

func InitLayoutCache(size int) {
	layoutCache.Lock()
	defer layoutCache.Unlock()

	layoutCache.size = maxInt(0, size)
	layoutCache.order = nil
	layoutCache.items = nil
	if layoutCache.size > 0 {
		layoutCache.items = make(map[layoutCacheKey]solvedLayout, layoutCache.size)
	}
}

func (l Layout) solveLayout(area Rect) solvedLayout {
	key := l.layoutCacheKey(area)
	if cached, ok := getCachedLayout(key); ok {
		return cached
	}

	solved := l.solveLayoutUncached(area)
	putCachedLayout(key, solved)
	return solved
}

func (l Layout) solveLayoutUncached(area Rect) solvedLayout {
	solver := casow.NewSolver()
	variables := make([]casow.Variable, len(l.constraints)*2+2)
	for i := range variables {
		variables[i] = casow.NewVariable()
	}

	spacers := make([]layoutElement, 0, len(l.constraints)+1)
	for i := 0; i+1 < len(variables); i += 2 {
		spacers = append(spacers, layoutElement{start: variables[i], end: variables[i+1]})
	}
	segments := make([]layoutElement, 0, len(l.constraints))
	for i := 1; i+1 < len(variables); i += 2 {
		segments = append(segments, layoutElement{start: variables[i], end: variables[i+1]})
	}

	areaElement := layoutElement{start: variables[0], end: variables[len(variables)-1]}
	areaStart, areaEnd := areaAxisBounds(area, l.direction)
	if err := configureArea(solver, areaElement, areaStart, areaEnd); err != nil {
		return solvedLayout{}
	}
	if err := configureVariableInAreaConstraints(solver, variables, areaElement); err != nil {
		return solvedLayout{}
	}
	if err := configureVariableConstraints(solver, variables); err != nil {
		return solvedLayout{}
	}
	if err := configureFlexConstraints(solver, areaElement, spacers, l.flex, l.spacing); err != nil {
		return solvedLayout{}
	}
	if err := configureConstraints(solver, areaElement, segments, l.constraints, l.flex); err != nil {
		return solvedLayout{}
	}
	if err := configureFillConstraints(solver, segments, l.constraints, l.flex); err != nil {
		return solvedLayout{}
	}
	if l.flex != FlexLegacy {
		for i := 0; i+1 < len(segments); i++ {
			if err := solver.AddConstraint(segments[i].hasSize(segments[i+1].size(), strengthAllSegmentGrow.casow())); err != nil {
				return solvedLayout{}
			}
		}
	}

	return solvedLayout{
		segments: elementsToRects(solver, segments, area, l.direction),
		spacers:  elementsToRects(solver, spacers, area, l.direction),
	}
}

type layoutElement struct {
	start casow.Variable
	end   casow.Variable
}

func (e layoutElement) size() casow.Expression {
	return casow.Var(e.end).MinusExpression(casow.Var(e.start))
}

func (e layoutElement) hasSize(size any, strength casow.Strength) casow.Constraint {
	return casow.NewConstraint(e.size(), casow.Equal, size, strength)
}

func (e layoutElement) hasMinSize(size int, strength casow.Strength) casow.Constraint {
	return casow.NewConstraint(e.size(), casow.GreaterOrEqual, float64(size), strength)
}

func (e layoutElement) hasMaxSize(size int, strength casow.Strength) casow.Constraint {
	return casow.NewConstraint(e.size(), casow.LessOrEqual, float64(size), strength)
}

func (e layoutElement) isEmpty() casow.Constraint {
	return e.hasSize(0, casow.Required.Sub(casow.Weak))
}

func areaAxisBounds(area Rect, direction Direction) (int, int) {
	if direction == Vertical {
		return area.Y, area.Y + area.Height
	}
	return area.X, area.X + area.Width
}

func configureArea(solver *casow.Solver, area layoutElement, areaStart int, areaEnd int) error {
	return solver.AddConstraints(
		casow.NewConstraint(area.start, casow.Equal, float64(areaStart), casow.Required),
		casow.NewConstraint(area.end, casow.Equal, float64(areaEnd), casow.Required),
	)
}

func configureVariableInAreaConstraints(solver *casow.Solver, variables []casow.Variable, area layoutElement) error {
	for _, variable := range variables {
		if err := solver.AddConstraints(
			casow.NewConstraint(variable, casow.GreaterOrEqual, area.start, casow.Required),
			casow.NewConstraint(variable, casow.LessOrEqual, area.end, casow.Required),
		); err != nil {
			return err
		}
	}
	return nil
}

func configureVariableConstraints(solver *casow.Solver, variables []casow.Variable) error {
	for i := 1; i+1 < len(variables); i += 2 {
		if err := solver.AddConstraint(casow.NewConstraint(variables[i], casow.LessOrEqual, variables[i+1], casow.Required)); err != nil {
			return err
		}
	}
	return nil
}

func configureFlexConstraints(solver *casow.Solver, area layoutElement, spacers []layoutElement, flex Flex, spacing int) error {
	middle := middleSpacers(spacers)
	switch flex {
	case FlexLegacy:
		if err := constrainEachSize(solver, middle, float64(spacing), strengthSpacerSize.casow()); err != nil {
			return err
		}
		return constrainOuterSpacersEmpty(solver, spacers)
	case FlexSpaceAround:
		if len(spacers) <= 2 {
			return configureEvenlyGrowingSpacers(solver, area, spacers, spacing)
		}
		first := spacers[0]
		last := spacers[len(spacers)-1]
		middle := spacers[1 : len(spacers)-1]
		if err := constrainAllEqual(solver, middle, strengthSpacerSize.casow()); err != nil {
			return err
		}
		if len(middle) > 0 {
			doubleFirst := casow.Var(first.end).MinusExpression(casow.Var(first.start)).Mul(2)
			doubleLast := casow.Var(last.end).MinusExpression(casow.Var(last.start)).Mul(2)
			if err := solver.AddConstraints(
				casow.NewConstraint(middle[0].size(), casow.Equal, doubleFirst, strengthSpacerSize.casow()),
				casow.NewConstraint(middle[0].size(), casow.Equal, doubleLast, strengthSpacerSize.casow()),
			); err != nil {
				return err
			}
		}
		return constrainGrowingSpacers(solver, area, spacers, spacing)
	case FlexSpaceEvenly:
		return configureEvenlyGrowingSpacers(solver, area, spacers, spacing)
	case FlexSpaceBetween:
		if err := constrainAllEqual(solver, middle, strengthSpacerSize.casow()); err != nil {
			return err
		}
		if err := constrainGrowingSpacers(solver, area, middle, spacing); err != nil {
			return err
		}
		return constrainOuterSpacersEmpty(solver, spacers)
	case FlexStart:
		if err := constrainEachSize(solver, middle, float64(spacing), strengthSpacerSize.casow()); err != nil {
			return err
		}
		return constrainStartCenterEnd(solver, area, spacers, true, false)
	case FlexCenter:
		if err := constrainEachSize(solver, middle, float64(spacing), strengthSpacerSize.casow()); err != nil {
			return err
		}
		return constrainStartCenterEnd(solver, area, spacers, false, false)
	case FlexEnd:
		if err := constrainEachSize(solver, middle, float64(spacing), strengthSpacerSize.casow()); err != nil {
			return err
		}
		return constrainStartCenterEnd(solver, area, spacers, false, true)
	default:
		if err := constrainEachSize(solver, middle, float64(spacing), strengthSpacerSize.casow()); err != nil {
			return err
		}
		return constrainStartCenterEnd(solver, area, spacers, true, false)
	}
}

func configureConstraints(solver *casow.Solver, area layoutElement, segments []layoutElement, constraints []Constraint, flex Flex) error {
	for i, constraint := range constraints {
		segment := segments[i]
		lastLegacyFlexible := flex == FlexLegacy && i == len(constraints)-1 && constraint.kind != constraintMin && constraint.kind != constraintFill && !hasEarlierGrowConstraint(constraints)
		maxSizeStrength := strengthMaxSize.casow()
		maxSizeEqStrength := strengthMaxSizeEq.casow()
		lengthSizeStrength := strengthLengthSize.casow()
		percentageSizeStrength := strengthPercentageSize.casow()
		ratioSizeStrength := strengthRatioSize.casow()
		if lastLegacyFlexible {
			maxSizeStrength = strengthAllSegmentGrow.casow()
			maxSizeEqStrength = strengthAllSegmentGrow.casow()
			lengthSizeStrength = strengthAllSegmentGrow.casow()
			percentageSizeStrength = strengthAllSegmentGrow.casow()
			ratioSizeStrength = strengthAllSegmentGrow.casow()
		}
		switch constraint.kind {
		case constraintMax:
			if err := solver.AddConstraints(
				segment.hasMaxSize(constraint.value, maxSizeStrength),
				segment.hasSize(float64(constraint.value), maxSizeEqStrength),
			); err != nil {
				return err
			}
		case constraintMin:
			if err := solver.AddConstraint(segment.hasMinSize(constraint.value, strengthMinSize.casow())); err != nil {
				return err
			}
			if flex == FlexLegacy {
				if err := solver.AddConstraint(segment.hasSize(float64(constraint.value), strengthMinSizeEq.casow())); err != nil {
					return err
				}
			} else if err := solver.AddConstraint(segment.hasSize(area.size(), strengthFillGrow.casow())); err != nil {
				return err
			}
		case constraintLength:
			if err := solver.AddConstraint(segment.hasSize(float64(constraint.value), lengthSizeStrength)); err != nil {
				return err
			}
		case constraintPercentage:
			size := area.size().Mul(float64(constraint.value)).Div(100)
			if err := solver.AddConstraint(segment.hasSize(size, percentageSizeStrength)); err != nil {
				return err
			}
		case constraintRatio:
			denominator := maxInt(1, constraint.denominator)
			size := area.size().Mul(float64(constraint.value)).Div(float64(denominator))
			if err := solver.AddConstraint(segment.hasSize(size, ratioSizeStrength)); err != nil {
				return err
			}
		case constraintFill:
			if err := solver.AddConstraint(segment.hasSize(area.size(), strengthFillGrow.casow())); err != nil {
				return err
			}
		}
	}
	return nil
}

func hasEarlierGrowConstraint(constraints []Constraint) bool {
	for i := 0; i+1 < len(constraints); i++ {
		if constraints[i].kind == constraintMin || constraints[i].kind == constraintFill {
			return true
		}
	}
	return false
}

func configureFillConstraints(solver *casow.Solver, segments []layoutElement, constraints []Constraint, flex Flex) error {
	for i := range constraints {
		if !isFillParticipant(constraints[i], flex) {
			continue
		}
		for j := i + 1; j < len(constraints); j++ {
			if !isFillParticipant(constraints[j], flex) {
				continue
			}
			leftScale := fillScale(constraints[i])
			rightScale := fillScale(constraints[j])
			left := segments[i].size().Mul(rightScale)
			right := segments[j].size().Mul(leftScale)
			if err := solver.AddConstraint(casow.NewConstraint(left, casow.Equal, right, strengthGrow.casow())); err != nil {
				return err
			}
		}
	}
	return nil
}

func isFillParticipant(constraint Constraint, flex Flex) bool {
	return constraint.kind == constraintFill || (flex != FlexLegacy && constraint.kind == constraintMin)
}

func fillScale(constraint Constraint) float64 {
	if constraint.kind == constraintFill {
		return math.Max(float64(constraint.value), 1e-6)
	}
	return 1
}

func configureEvenlyGrowingSpacers(solver *casow.Solver, area layoutElement, spacers []layoutElement, spacing int) error {
	if err := constrainAllEqual(solver, spacers, strengthSpacerSize.casow()); err != nil {
		return err
	}
	return constrainGrowingSpacers(solver, area, spacers, spacing)
}

func constrainGrowingSpacers(solver *casow.Solver, area layoutElement, spacers []layoutElement, spacing int) error {
	for _, spacer := range spacers {
		if err := solver.AddConstraints(
			spacer.hasMinSize(spacing, strengthSpacerSize.casow()),
			spacer.hasSize(area.size(), strengthSpaceGrow.casow()),
		); err != nil {
			return err
		}
	}
	return nil
}

func constrainAllEqual(solver *casow.Solver, elements []layoutElement, strength casow.Strength) error {
	for i := range elements {
		for j := i + 1; j < len(elements); j++ {
			if err := solver.AddConstraint(elements[i].hasSize(elements[j].size(), strength)); err != nil {
				return err
			}
		}
	}
	return nil
}

func constrainEachSize(solver *casow.Solver, elements []layoutElement, size float64, strength casow.Strength) error {
	for _, element := range elements {
		if err := solver.AddConstraint(element.hasSize(size, strength)); err != nil {
			return err
		}
	}
	return nil
}

func constrainOuterSpacersEmpty(solver *casow.Solver, spacers []layoutElement) error {
	if len(spacers) == 0 {
		return nil
	}
	if err := solver.AddConstraint(spacers[0].isEmpty()); err != nil {
		return err
	}
	return solver.AddConstraint(spacers[len(spacers)-1].isEmpty())
}

func constrainStartCenterEnd(solver *casow.Solver, area layoutElement, spacers []layoutElement, emptyFirst bool, emptyLast bool) error {
	if len(spacers) == 0 {
		return nil
	}
	first := spacers[0]
	last := spacers[len(spacers)-1]
	if emptyFirst {
		if err := solver.AddConstraint(first.isEmpty()); err != nil {
			return err
		}
		if err := solver.AddConstraint(last.hasSize(area.size(), strengthGrow.casow())); err != nil {
			return err
		}
		return nil
	}
	if emptyLast {
		if err := solver.AddConstraint(last.isEmpty()); err != nil {
			return err
		}
		if err := solver.AddConstraint(first.hasSize(area.size(), strengthGrow.casow())); err != nil {
			return err
		}
		return nil
	}
	return solver.AddConstraints(
		first.hasSize(area.size(), strengthGrow.casow()),
		last.hasSize(area.size(), strengthGrow.casow()),
		first.hasSize(last.size(), strengthSpacerSize.casow()),
	)
}

func middleSpacers(spacers []layoutElement) []layoutElement {
	if len(spacers) <= 2 {
		return nil
	}
	return spacers[1 : len(spacers)-1]
}

func elementsToRects(solver *casow.Solver, elements []layoutElement, area Rect, direction Direction) []Rect {
	rects := make([]Rect, len(elements))
	for i, element := range elements {
		start := roundSolvedValue(solver.GetValue(element.start))
		end := roundSolvedValue(solver.GetValue(element.end))
		size := maxInt(0, end-start)
		if direction == Horizontal {
			rects[i] = Rect{X: start, Y: area.Y, Width: size, Height: area.Height}
		} else {
			rects[i] = Rect{X: area.X, Y: start, Width: area.Width, Height: size}
		}
	}
	return rects
}

func roundSolvedValue(value float64) int {
	if value <= 0 {
		return 0
	}
	return int(math.Round(value + 1e-9))
}

func (l Layout) layoutCacheKey(area Rect) layoutCacheKey {
	return layoutCacheKey{
		area:        area,
		direction:   l.direction,
		flex:        l.flex,
		spacing:     l.spacing,
		margin:      l.margin,
		constraints: fmt.Sprint(l.constraints),
	}
}

func getCachedLayout(key layoutCacheKey) (solvedLayout, bool) {
	layoutCache.Lock()
	defer layoutCache.Unlock()

	if layoutCache.size <= 0 || layoutCache.items == nil {
		return solvedLayout{}, false
	}
	value, ok := layoutCache.items[key]
	if !ok {
		return solvedLayout{}, false
	}
	return cloneSolvedLayout(value), true
}

func putCachedLayout(key layoutCacheKey, value solvedLayout) {
	layoutCache.Lock()
	defer layoutCache.Unlock()

	if layoutCache.size <= 0 {
		return
	}
	if layoutCache.items == nil {
		layoutCache.items = make(map[layoutCacheKey]solvedLayout, layoutCache.size)
	}
	if _, ok := layoutCache.items[key]; !ok {
		layoutCache.order = append(layoutCache.order, key)
	}
	layoutCache.items[key] = cloneSolvedLayout(value)

	for len(layoutCache.order) > layoutCache.size {
		oldest := layoutCache.order[0]
		layoutCache.order = layoutCache.order[1:]
		delete(layoutCache.items, oldest)
	}
}

func cloneSolvedLayout(value solvedLayout) solvedLayout {
	return solvedLayout{
		segments: append([]Rect(nil), value.segments...),
		spacers:  append([]Rect(nil), value.spacers...),
	}
}
