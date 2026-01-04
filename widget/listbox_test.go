package widget

import "testing"

func TestListBoxCursorClamping_REQ_UI_PRIMITIVE_TESTS(t *testing.T) {
	// [REQ:UI_PRIMITIVE_TESTS] [ARCH:TEST_STRATEGY_UI] [IMPL:TEST_WIDGETS]
	// Validates ListBox cursor math and fallback semantics.
	lb := NewListBox(0, 0, 10, 5, "cursor")
	lb.SetLower(2)
	lb.AppendString("alpha", "bravo", "charlie", "delta", "echo")

	lb.SetCursor(-5)
	if lb.Cursor() != lb.Lower() {
		t.Fatalf("SetCursor below lower bound -> %d, want %d", lb.Cursor(), lb.Lower())
	}

	lb.SetCursor(999)
	want := lb.Upper() - 1
	if lb.Cursor() != want {
		t.Fatalf("SetCursor above upper bound -> %d, want %d", lb.Cursor(), want)
	}

	lb.SetCursorByName("charlie")
	if lb.Cursor() != 2 {
		t.Fatalf("SetCursorByName charlie -> %d, want 2", lb.Cursor())
	}

	lb.SetCursorByName("missing")
	if lb.Cursor() != lb.Lower() {
		t.Fatalf("SetCursorByName fallback -> %d, want lower %d", lb.Cursor(), lb.Lower())
	}
}

func TestListBoxScrollRate_REQ_UI_PRIMITIVE_TESTS(t *testing.T) {
	// [REQ:UI_PRIMITIVE_TESTS] [ARCH:TEST_STRATEGY_UI] [IMPL:TEST_WIDGETS]
	lb := NewListBox(0, 0, 12, 6, "scroll")
	for i := 0; i < 12; i++ {
		lb.AppendString(t.Name() + string(rune('A'+i)))
	}

	if got := lb.ScrollRate(); got != "Top" {
		t.Fatalf("ScrollRate at offset 0 -> %s, want Top", got)
	}

	lb.offset = lb.rowCol() // mid point (rowCol = (6-2)=4)
	if got := lb.ScrollRate(); got != "50%" {
		t.Fatalf("ScrollRate mid offset -> %s, want 50%%", got)
	}

	lb.offset = lb.Upper() - lb.rowCol()
	if got := lb.ScrollRate(); got != "Bot" {
		t.Fatalf("ScrollRate bottom -> %s, want Bot", got)
	}
}

func TestListBoxColumnAdjust_REQ_UI_PRIMITIVE_TESTS(t *testing.T) {
	// [REQ:UI_PRIMITIVE_TESTS] [ARCH:TEST_STRATEGY_UI] [IMPL:TEST_WIDGETS]
	lb := NewListBox(0, 0, 24, 5, "columns")
	lb.AppendString("aaaa", "bbbbbb", "1234567890")

	lb.ColumnAdjustContentsWidth()
	if lb.Column() != 2 {
		t.Fatalf("ColumnAdjustContentsWidth -> %d, want 2", lb.Column())
	}
}
