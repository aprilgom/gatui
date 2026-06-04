package text_test

import (
	"fmt"
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/text"
)

func TestMasked_New_shouldStoreMaskChar(t *testing.T) {
	masked := text.NewMasked("12345", 'x')

	if got := masked.MaskChar(); got != 'x' {
		t.Fatalf("MaskChar() = %q, want %q", got, 'x')
	}
}

func TestMasked_Value_shouldReturnMaskedString(t *testing.T) {
	masked := text.NewMasked("12345", 'x')

	if got := masked.Value(); got != "xxxxx" {
		t.Fatalf("Value() = %q, want %q", got, "xxxxx")
	}
}

func TestMasked_String_shouldDisplayMaskedString(t *testing.T) {
	masked := text.NewMasked("12345", 'x')

	if got := fmt.Sprint(masked); got != "xxxxx" {
		t.Fatalf("String() = %q, want %q", got, "xxxxx")
	}
}

func TestMasked_Text_shouldConvertToText(t *testing.T) {
	masked := text.NewMasked("12345", 'x')

	got := masked.Text()

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	if len(got.Lines[0].Spans) != 1 {
		t.Fatalf("span count = %d, want 1", len(got.Lines[0].Spans))
	}
	if got.Lines[0].Spans[0].Content != "xxxxx" {
		t.Fatalf("content = %q, want %q", got.Lines[0].Spans[0].Content, "xxxxx")
	}
}

func TestMasked_Text_shouldRenderMaskedValue(t *testing.T) {
	masked := text.NewMasked("12345", 'x')
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	masked.Text().Render(buf.Area, buf)

	if got := buf.Lines(); len(got) != 1 || got[0] != "xxxxx" {
		t.Fatalf("rendered lines = %#v, want %#v", got, []string{"xxxxx"})
	}
}

func TestMasked_Value_shouldCountRunesNotBytes(t *testing.T) {
	masked := text.NewMasked("가나다", '*')

	if got := masked.Value(); got != "***" {
		t.Fatalf("Value() = %q, want %q", got, "***")
	}
}

func TestMasked_Value_shouldSupportWideMaskRune(t *testing.T) {
	masked := text.NewMasked("abc", '＊')

	if got := masked.Value(); got != "＊＊＊" {
		t.Fatalf("Value() = %q, want %q", got, "＊＊＊")
	}
	if got := masked.Text().Width(); got != 6 {
		t.Fatalf("Text().Width() = %d, want 6", got)
	}
}
