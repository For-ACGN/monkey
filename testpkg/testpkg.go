package testpkg

import (
	"fmt"
	"io"
)

// Writer for test.
type Writer struct {
	data []byte
}

// implement io.Writer.
func (w *Writer) Write(b []byte) (int, error) {
	w.data = b
	return w.print()
}

// Print padding comment.
func (w *Writer) Print() (int, error) {
	n, err := w.print()
	fmt.Println()
	return n + 1, err
}

// Println padding comment.
func (w Writer) Println() (int, error) {
	n, err := w.println()
	fmt.Println()
	return n, err
}

func (w *Writer) print() (int, error) {
	return fmt.Print(string(w.data) + "\n")
}

func (w Writer) println() (int, error) {
	return fmt.Println(string(w.data))
}

// writer for test.
type writer struct {
	data []byte
}

// NewWriter is used to create a private writer.
func NewWriter() io.Writer {
	return new(writer)
}

// implement io.Writer.
func (w *writer) Write(b []byte) (int, error) {
	w.data = b
	return w.print()
}

func (w *writer) print() (int, error) {
	return fmt.Print(string(w.data) + "\n")
}
