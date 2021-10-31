package main

import (
	"regexp"
	"strings"
)

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

// HeadingRE regex source: https://github.com/Python-Markdown/markdown/blob/master/markdown/blockprocessors.py#L448
const HeadingRE = `(?:^|\n)(?P<level>#{1,6})(?P<header>(?:\\.|[^\\])*?)#*(?:\n|$)`

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
	p := &mdParser{}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// headings
		r := regexp.MustCompile(HeadingRE)
		if r.MatchString(line) {
			level := len(r.FindStringSubmatch(line)[1])
			content := strings.TrimSpace(r.FindStringSubmatch(line)[2])
			/*
			* heading 1-6 has value 1-6 according to the enum declaration above.
			* So, we can directly use the length of `#` to determine the heading level
			 */
			p.tokens = append(p.tokens, &token{level, content})
			continue
		}

		// paragraph
		p.tokens = append(p.tokens, &token{para, line})
	}

	return p
}
