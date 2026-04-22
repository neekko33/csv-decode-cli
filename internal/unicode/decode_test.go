package unicode

import "testing"

func TestDecodeEscapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "u escape", input: `\u65e5\u672c`, want: "日本"},
		{name: "U escape", input: `\U000030c6\U000030b9\U000030c8`, want: "テスト"},
		{name: "mixed text", input: `abc-\u4eca\u65e5`, want: "abc-今日"},
		{name: "invalid escape preserved", input: `\uZZZZ`, want: `\uZZZZ`},
		{name: "regular text", input: "plain", want: "plain"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DecodeEscapes(tc.input)
			if got != tc.want {
				t.Fatalf("DecodeEscapes() = %q, want %q", got, tc.want)
			}
		})
	}
}
