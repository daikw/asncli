package cli

import (
	"fmt"
	"io"
	"strings"
)

type Renderer interface {
	JSON(data any) error
	Envelope(data any, next any, warnings []string) error
	Table(headers []string, rows [][]string) error
	Message(format string, args ...any) error
}

type DefaultRenderer struct {
	output Output
}

func NewRenderer(stdout, stderr io.Writer, jsonOutput bool) Renderer {
	return DefaultRenderer{output: Output{Stdout: stdout, Stderr: stderr, JSON: jsonOutput}}
}

func (r DefaultRenderer) JSON(data any) error {
	return r.output.PrintJSON(data)
}

func (r DefaultRenderer) Envelope(data any, next any, warnings []string) error {
	return r.output.PrintEnvelope(Envelope{Data: data, Next: next, Warnings: warnings})
}

func (r DefaultRenderer) Table(headers []string, rows [][]string) error {
	w := r.output.Table()
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return w.Flush()
}

func (r DefaultRenderer) Message(format string, args ...any) error {
	_, err := fmt.Fprintf(r.output.Stdout, format, args...)
	return err
}
