package cli

import (
	"encoding/json"
	"io"
	"text/tabwriter"
)

type Envelope struct {
	Data     any      `json:"data"`
	Next     any      `json:"next,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type Output struct {
	Stdout io.Writer
	Stderr io.Writer
	JSON   bool
}

func (o Output) PrintJSON(data any) error {
	enc := json.NewEncoder(o.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(Envelope{Data: data})
}

func (o Output) PrintEnvelope(env Envelope) error {
	enc := json.NewEncoder(o.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}

func (o Output) Table() *tabwriter.Writer {
	return tabwriter.NewWriter(o.Stdout, 0, 4, 2, ' ', 0)
}
