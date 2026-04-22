package csvsvc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultOutputPath(t *testing.T) {
	t.Parallel()

	in := filepath.Join("tmp", "input.csv")
	got := DefaultOutputPath(in)
	want := filepath.Join("tmp", "input-decoded.csv")
	if got != want {
		t.Fatalf("DefaultOutputPath() = %q, want %q", got, want)
	}
}

func TestValidateDestination(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	existing := filepath.Join(tmp, "exists.csv")
	if err := os.WriteFile(existing, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateDestination("", false); err == nil {
		t.Fatal("expected error for empty output path")
	}

	err := ValidateDestination(existing, false)
	if !errors.Is(err, ErrDestinationExists) {
		t.Fatalf("expected ErrDestinationExists, got %v", err)
	}

	if err := ValidateDestination(existing, true); err != nil {
		t.Fatalf("expected overwrite allowed, got %v", err)
	}
}

func TestDecodeCSVFields(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	in := filepath.Join(tmp, "input.csv")
	out := filepath.Join(tmp, "output.csv")

	content := "id,message,title\n1,\\u65e5\\u672c,\\u30c6\\u30b9\\u30c8\n2,plain,\\u4eca\\u65e5\n"
	if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := DecodeCSVFields(in, out, []string{"message", "title"}, false); err != nil {
		t.Fatalf("DecodeCSVFields() error = %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	got := string(data)
	if !strings.Contains(got, "日本") || !strings.Contains(got, "テスト") || !strings.Contains(got, "今日") {
		t.Fatalf("decoded CSV missing expected text:\n%s", got)
	}
}

func TestDecodeCSVFields_FieldNotFound(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	in := filepath.Join(tmp, "input.csv")
	out := filepath.Join(tmp, "output.csv")

	if err := os.WriteFile(in, []byte("id,name\n1,foo\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := DecodeCSVFields(in, out, []string{"message"}, false)
	if err == nil || !strings.Contains(err.Error(), `field "message" not found`) {
		t.Fatalf("expected field not found error, got %v", err)
	}
}
