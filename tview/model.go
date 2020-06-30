package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Option represents a possible choice at the end of an arc.
type Option struct {
	Label string `json:"text"`
	URL   string `json:"arc"`
}

// Arc represents an individual narrative structure.
type Arc struct {
	Title   string   `json:"title"`
	Body    []string `json:"story"`
	Options []Option `json:"options"`
}

// Story is the graph of arcs.
// No Mutex: we never write after read is done.
type Story map[string]Arc

// Load fetches the story data from disk.
func (s *Story) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`opening %s: %w`, path, err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(s)
	return err
}

// Arc obtains an arc from its URL, returning nil if no arc is found for that URL.
func (s *Story) Arc(url string) *Arc {
	a, ok := (*s)[url]
	if !ok {
		return nil
	}
	return &a
}
