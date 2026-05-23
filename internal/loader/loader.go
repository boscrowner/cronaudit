// Package loader provides utilities for reading crontab entries
// from files or standard input.
package loader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Source represents the origin of crontab input.
type Source struct {
	Name   string
	Reader io.Reader
}

// FromFile opens the given file path and returns a Source for reading.
func FromFile(path string) (*Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: cannot open file %q: %w", path, err)
	}
	return &Source{Name: path, Reader: f}, nil
}

// FromStdin returns a Source that reads from standard input.
func FromStdin() *Source {
	return &Source{Name: "<stdin>", Reader: os.Stdin}
}

// Lines reads all lines from the Source and returns them as a slice of strings.
// Empty lines and lines beginning with '#' are included as-is; the parser
// is responsible for filtering them.
func Lines(s *Source) ([]string, error) {
	if s == nil {
		return nil, fmt.Errorf("loader: nil source")
	}

	var lines []string
	scanner := bufio.NewScanner(s.Reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: error reading %q: %w", s.Name, err)
	}
	return lines, nil
}

// Close attempts to close the underlying reader if it implements io.Closer.
func Close(s *Source) error {
	if s == nil {
		return nil
	}
	if c, ok := s.Reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
