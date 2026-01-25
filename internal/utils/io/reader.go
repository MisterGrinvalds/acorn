package io

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Reader handles input deserialization with streaming support.
type Reader struct {
	config      *IOConfig
	reader      io.Reader
	bufReader   *bufio.Reader
	inputMode   InputMode
	initialized bool
	file        *os.File // If we opened a file, track it for closing
}

// NewReader creates a Reader from IOConfig.
func NewReader(cfg *IOConfig) (*Reader, error) {
	r := &Reader{
		config: cfg,
	}

	// Determine input source
	if cfg.InputFile != "" {
		f, err := os.Open(cfg.InputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open input file: %w", err)
		}
		r.reader = f
		r.file = f
	} else if cfg.InputReader != nil {
		r.reader = cfg.InputReader
	} else {
		r.reader = os.Stdin
	}

	return r, nil
}

// initialize sets up the buffered reader and detects format if needed.
func (r *Reader) initialize() error {
	if r.initialized {
		return nil
	}

	// Detect format if auto
	if r.config.InputFormat == FormatAuto {
		mode, newReader, err := DetectInputMode(r.reader)
		if err != nil {
			return fmt.Errorf("failed to detect input format: %w", err)
		}
		r.inputMode = mode
		r.reader = newReader
	} else {
		// Convert Format to InputMode
		switch r.config.InputFormat {
		case FormatJSON:
			r.inputMode = ModeBatchJSON
		case FormatNDJSON:
			r.inputMode = ModeNDJSON
		case FormatYAML:
			r.inputMode = ModeBatchYAML
		default:
			r.inputMode = ModeBatchJSON
		}
	}

	// Create buffered reader if not already buffered
	if br, ok := r.reader.(*bufio.Reader); ok {
		r.bufReader = br
	} else {
		r.bufReader = bufio.NewReaderSize(r.reader, 64*1024)
	}

	r.initialized = true
	return nil
}

// Read deserializes the entire input into the target type.
func (r *Reader) Read(target interface{}) error {
	if err := r.initialize(); err != nil {
		return err
	}

	data, err := io.ReadAll(r.bufReader)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if len(data) == 0 {
		return nil // Empty input is valid
	}

	return r.unmarshal(data, target)
}

// ReadRaw reads the input as raw bytes without deserialization.
func (r *Reader) ReadRaw() ([]byte, error) {
	if err := r.initialize(); err != nil {
		return nil, err
	}

	return io.ReadAll(r.bufReader)
}

// StreamRecord represents a single record from streaming input.
type StreamRecord struct {
	Data    interface{}
	Index   int
	LineNum int
	Raw     []byte
	Error   error
}

// ReadStream returns a channel that yields items one at a time.
// For NDJSON, each line is a separate item.
// For JSON arrays, each array element is a separate item.
// For YAML, each document is a separate item.
func (r *Reader) ReadStream() (<-chan *StreamRecord, error) {
	if err := r.initialize(); err != nil {
		return nil, err
	}

	ch := make(chan *StreamRecord, 100)

	go func() {
		defer close(ch)

		switch r.inputMode {
		case ModeNDJSON:
			r.streamNDJSON(ch)
		case ModeBatchJSON:
			r.streamJSONArray(ch)
		case ModeBatchYAML, ModeMultiYAML:
			r.streamYAML(ch)
		default:
			// Read all and emit as single item
			var data interface{}
			if err := r.Read(&data); err != nil {
				ch <- &StreamRecord{Error: err}
				return
			}
			ch <- &StreamRecord{Data: data, Index: 0}
		}
	}()

	return ch, nil
}

// streamNDJSON reads JSON Lines format (one JSON object per line).
func (r *Reader) streamNDJSON(ch chan<- *StreamRecord) {
	scanner := bufio.NewScanner(r.bufReader)
	// Handle large JSON lines
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

	lineNum := 0
	index := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(line) == 0 {
			continue // Skip empty lines
		}

		var obj interface{}
		if err := json.Unmarshal(line, &obj); err != nil {
			ch <- &StreamRecord{
				Error:   fmt.Errorf("line %d: %w", lineNum, err),
				LineNum: lineNum,
				Raw:     append([]byte{}, line...),
			}
			continue
		}

		ch <- &StreamRecord{
			Data:    obj,
			Index:   index,
			LineNum: lineNum,
			Raw:     append([]byte{}, line...),
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		ch <- &StreamRecord{Error: fmt.Errorf("scanner error: %w", err)}
	}
}

// streamJSONArray reads a JSON array and emits each element.
func (r *Reader) streamJSONArray(ch chan<- *StreamRecord) {
	dec := json.NewDecoder(r.bufReader)
	dec.UseNumber()

	// Try to read opening bracket
	tok, err := dec.Token()
	if err != nil {
		if err == io.EOF {
			return // Empty input
		}
		ch <- &StreamRecord{Error: fmt.Errorf("failed to read JSON: %w", err)}
		return
	}

	// Check if it's an array
	if delim, ok := tok.(json.Delim); ok && delim == '[' {
		// It's an array, read elements
		index := 0
		for dec.More() {
			var obj interface{}
			if err := dec.Decode(&obj); err != nil {
				ch <- &StreamRecord{Error: fmt.Errorf("element %d: %w", index, err)}
				return
			}
			ch <- &StreamRecord{Data: obj, Index: index}
			index++
		}
	} else {
		// It's a single object, emit as one item
		// We already consumed the first token, so we need to handle this differently
		// Read the rest and combine with the first token
		data, err := io.ReadAll(dec.Buffered())
		if err != nil {
			ch <- &StreamRecord{Error: fmt.Errorf("failed to read remaining data: %w", err)}
			return
		}
		rest, _ := io.ReadAll(r.bufReader)
		data = append(data, rest...)

		// Reconstruct the full JSON
		var fullData []byte
		switch v := tok.(type) {
		case string:
			fullData = append([]byte(`"`+v+`"`), data...)
		case json.Number:
			fullData = append([]byte(v.String()), data...)
		case bool:
			if v {
				fullData = append([]byte("true"), data...)
			} else {
				fullData = append([]byte("false"), data...)
			}
		case nil:
			fullData = append([]byte("null"), data...)
		default:
			fullData = data
		}

		var obj interface{}
		if err := json.Unmarshal(fullData, &obj); err != nil {
			ch <- &StreamRecord{Error: fmt.Errorf("failed to parse JSON object: %w", err)}
			return
		}
		ch <- &StreamRecord{Data: obj, Index: 0}
	}
}

// streamYAML reads YAML documents and emits each one.
func (r *Reader) streamYAML(ch chan<- *StreamRecord) {
	dec := yaml.NewDecoder(r.bufReader)
	index := 0

	for {
		var doc interface{}
		err := dec.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			ch <- &StreamRecord{Error: fmt.Errorf("document %d: %w", index, err)}
			return
		}

		ch <- &StreamRecord{Data: doc, Index: index}
		index++
	}
}

// unmarshal deserializes data based on the detected/configured format.
func (r *Reader) unmarshal(data []byte, target interface{}) error {
	switch r.inputMode {
	case ModeBatchJSON, ModeNDJSON:
		return json.Unmarshal(data, target)
	case ModeBatchYAML, ModeMultiYAML:
		return yaml.Unmarshal(data, target)
	default:
		// Try JSON first, then YAML
		if err := json.Unmarshal(data, target); err == nil {
			return nil
		}
		return yaml.Unmarshal(data, target)
	}
}

// Close closes the underlying reader if it's a file.
func (r *Reader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// HasInput returns true if there's input available.
// For stdin, checks if it's a pipe or has content.
func (r *Reader) HasInput() bool {
	if r.config.InputFile != "" {
		return true
	}
	if r.config.InputReader != nil {
		return true
	}
	return HasStdinData()
}

// Format returns the detected or configured input format.
func (r *Reader) Format() Format {
	if r.inputMode != "" {
		return FormatFromInputMode(r.inputMode)
	}
	return r.config.InputFormat
}
