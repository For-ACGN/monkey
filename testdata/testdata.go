package testdata

import (
	"fmt"
)

type Writer struct {
	data []byte
}

func (w *Writer) Write(b []byte) (int, error) {
	w.data = b
	return w.print()
}

func (w *Writer) Print() (int, error) {
	n, err := w.print()
	fmt.Println()
	return n + 1, err
}

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
