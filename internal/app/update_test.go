package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error = %v", err)
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "tilde only",
			in:   "~",
			want: home,
		},
		{
			name: "tilde path",
			in:   filepath.Join("~", "Documents", "file.csv"),
			want: filepath.Join(home, "Documents", "file.csv"),
		},
		{
			name: "non tilde path",
			in:   "/tmp/out.csv",
			want: "/tmp/out.csv",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := setHomeDir(tc.in)
			if got != tc.want {
				t.Fatalf("setHomeDir(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
