package configfile

import (
	"testing"
)

func TestRawWriterFormat(t *testing.T) {
	w := &RawWriter{}
	if w.Format() != "raw" {
		t.Errorf("Format() = %q, want %q", w.Format(), "raw")
	}
}

func TestRawWriterContent(t *testing.T) {
	w := &RawWriter{}

	values := map[string]any{
		"content": "Hello, World!\nThis is raw text.",
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	expected := "Hello, World!\nThis is raw text."
	if string(content) != expected {
		t.Errorf("Content = %q, want %q", string(content), expected)
	}
}

func TestRawWriterMultiline(t *testing.T) {
	w := &RawWriter{}

	script := `.First <- function(){
    cat("\n   Loading Site-wide Profile\n\n")
}
.Last <- function()  cat("\n   Work done!\n\n")
`

	values := map[string]any{
		"content": script,
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if string(content) != script {
		t.Errorf("Content mismatch:\ngot:\n%s\nwant:\n%s", string(content), script)
	}
}

func TestRawWriterEmpty(t *testing.T) {
	w := &RawWriter{}

	values := map[string]any{}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty content, got: %s", string(content))
	}
}
