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
	strengthAreaBounds    = solverStrength{value: 1_001_001_000}
	strengthSegmentSizing = solverStrength{value: 1_000_000}
	strengthFlexSpacing   = solverStrength{value: 1_000}
	strengthSpacerGrowth  = solverStrength{value: 1}
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
	rects, offsets, lengths := l.splitSegments(area)
	rects = solveSegmentRects(area, l.direction, offsets, lengths, rects)
	return solvedLayout{segments: rects, spacers: l.spacerRects(area, offsets, lengths)}
}

func solveSegmentRects(area Rect, direction Direction, offsets []int, lengths []int, fallback []Rect) []Rect {
	solver := casow.NewSolver()
	starts := make([]casow.Variable, len(lengths))
	ends := make([]casow.Variable, len(lengths))
	solvedLengths := make([]casow.Variable, len(lengths))
	areaLength := axisLength(area, direction)

	for i := range lengths {
		starts[i] = casow.NewVariable()
		ends[i] = casow.NewVariable()
		solvedLengths[i] = casow.NewVariable()
		if err := addSegmentConstraints(solver, starts[i], ends[i], solvedLengths[i], offsets[i], lengths[i], areaLength); err != nil {
			return fallback
		}
	}

	rects := make([]Rect, len(fallback))
	for i, rect := range fallback {
		start := roundSolvedValue(solver.GetValue(starts[i]))
		length := roundSolvedValue(solver.GetValue(solvedLengths[i]))
		if direction == Horizontal {
			rect.X = area.X + start
			rect.Width = length
		} else {
			rect.Y = area.Y + start
			rect.Height = length
		}
		rects[i] = rect
	}
	return rects
}

func addSegmentConstraints(solver *casow.Solver, start, end, length casow.Variable, offset, segmentLength, areaLength int) error {
	constraints := []casow.Constraint{
		casow.NewConstraint(start, casow.GreaterOrEqual, 0, strengthAreaBounds.casow()),
		casow.NewConstraint(end, casow.GreaterOrEqual, start, strengthAreaBounds.casow()),
		casow.NewConstraint(end, casow.LessOrEqual, float64(areaLength), strengthAreaBounds.casow()),
		casow.NewConstraint(casow.Var(start).PlusExpression(casow.Var(length)), casow.Equal, end, strengthAreaBounds.casow()),
		casow.NewConstraint(start, casow.Equal, float64(offset), strengthFlexSpacing.casow()),
		casow.NewConstraint(length, casow.Equal, float64(segmentLength), strengthSegmentSizing.casow()),
	}
	return solver.AddConstraints(constraints...)
}

func roundSolvedValue(value float64) int {
	if value <= 0 {
		return 0
	}
	return int(math.Round(value))
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
