package main

import (
	"strings"
)

type Wordcount string

func (w Wordcount) Words() int {
	return len(strings.Fields(string(w)))
}

func (w Wordcount) Lines() int {
	return len(strings.Split(string(w), "\n"))
}

func (w Wordcount) Characters() int {
	runes := []rune(w)
	return len(runes)
}
