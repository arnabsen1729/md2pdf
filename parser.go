package main

import "strings"

/*
	There is no default structure in Golang called enums (enumerators).
	We can simulate it using constant variable declaration.
	Link: https://www.sohamkamani.com/golang/enums/
*/

const (
	para int = iota
	heading1
	heading2
	heading3
	heading4
	heading5
	heading6
)

type token struct {
	style   int
	content string
}

// mdParser will parse the text passed into markdown tokens.
type mdParser struct {
	tokens []*token
}

// newParser returns a pointer to an object of mdParser.
// It parses the passed string into markdown tokens.
func newParser(input string) *mdParser {
	input = strings.TrimSpace(input)
	p := &mdParser{}

	if strings.HasPrefix(input, "##") {
		t := &token{}
		t.content = strings.TrimLeftFunc(input, func(r rune) bool { return r == '#' || r == ' ' })
		t.style = heading2
		p.tokens = append(p.tokens, t)
	} else if strings.HasPrefix(input, "#") {
		t := &token{}
		t.content = strings.TrimLeftFunc(input, func(r rune) bool { return r == '#' || r == ' ' })
		t.style = heading1
		p.tokens = append(p.tokens, t)
	} else {
		t := &token{}
		t.content = input
		p.tokens = append(p.tokens, t)
	}

	return p
}
