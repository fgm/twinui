package main

import (
	"encoding/json"
	"os"
)

// Option represents a possible choice at the end of an arc.
type Option struct {
	Text    string
	ArcName string `json:"arc"`
}

// Arc represents an individuall narrative structure.
type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

// Story is the graph of arcs.
type Story map[string]Arc

func loadStory(path string) (Story, error) {
	story := Story{}
	file, err := os.Open(path)
	if err != nil {
		return story, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&story)
	return story, err
}
