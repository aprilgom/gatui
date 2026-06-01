package widgets

import (
	"gatui/buffer"
	"gatui/layout"
)

type Widget interface {
	Render(area layout.Rect, buf *buffer.Buffer)
}

type WidgetRef interface {
	RenderRef(area layout.Rect, buf *buffer.Buffer)
}

type StatefulWidget interface {
	RenderStatefulRef(area layout.Rect, buf *buffer.Buffer, state any)
}
