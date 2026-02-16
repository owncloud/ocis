package conversions

import "testing"

var scenarios = []struct {
	name      string
	input     string
	separator string
	out       []string
}{
	{
		"comma separated input",
		"a, b, c, d",
		",",
		[]string{"a", "b", "c", "d"},
	}, {
		"space separated input",
		"a b c d",
		" ",
		[]string{"a", "b", "c", "d"},
	},
}

func TestStringToSliceString(t *testing.T) {
	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			s := StringToSliceString(tt.input, tt.separator)
			for i, v := range tt.out {
				if s[i] != v {
					t.Errorf("got %q, want %q", s, tt.out)
				}
			}
		})
	}
}
