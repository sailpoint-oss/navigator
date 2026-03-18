package navigator

import "testing"

func TestContainsPosition(t *testing.T) {
	r := Range{
		Start: Position{Line: 2, Character: 5},
		End:   Position{Line: 4, Character: 10},
	}

	tests := []struct {
		pos  Position
		want bool
	}{
		{Position{Line: 3, Character: 0}, true},
		{Position{Line: 2, Character: 5}, true},
		{Position{Line: 4, Character: 9}, true},
		{Position{Line: 4, Character: 10}, false},
		{Position{Line: 1, Character: 0}, false},
		{Position{Line: 5, Character: 0}, false},
		{Position{Line: 2, Character: 4}, false},
	}
	for _, tt := range tests {
		got := ContainsPosition(r, tt.pos)
		if got != tt.want {
			t.Errorf("ContainsPosition(r, %v) = %v, want %v", tt.pos, got, tt.want)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	if !IsEmpty(Range{}) {
		t.Error("zero range should be empty")
	}
	if IsEmpty(Range{Start: Position{0, 0}, End: Position{0, 5}}) {
		t.Error("non-zero range should not be empty")
	}
}
