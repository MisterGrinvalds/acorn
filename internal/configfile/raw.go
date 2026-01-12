package configfile

// RawWriter implements the Writer interface for raw text content.
// Used for scripts, plain text files, or any content that should be output as-is.
type RawWriter struct{}

func init() {
	Register(&RawWriter{})
}

// Format returns the format identifier.
func (w *RawWriter) Format() string {
	return "raw"
}

// Write outputs the content value as raw text.
// Expected structure:
//
//	content: |
//	  # This is raw text content
//	  It will be output exactly as written
func (w *RawWriter) Write(values map[string]interface{}) ([]byte, error) {
	content, ok := values["content"].(string)
	if !ok {
		return []byte{}, nil
	}
	return []byte(content), nil
}
