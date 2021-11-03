package main

import (
	"fmt"
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
	bold
	italic
)

// HeadingRE regex source: https://github.com/Python-Markdown/markdown/blob/master/markdown/blockprocessors.py#L448
const (
	HeadingRE = `(?:^|\n)(?P<level>#{1,6})(?P<header>(?:\\.|[^\\])*?)#*(?:\n|$)`
	BoldRE    = `(\*{2}\S+?\*{2})`
	ItalicRE  = `(\*{1}\S+?\*{1})`
)

type token struct {
	style   int
	content string
}

// mdParser will parse the text passed into markdown tokens.
type mdParser struct {
	lines [][]*token
}

func newTokenDerived(content string, style int) *token {
	content = strings.TrimSpace(content)

	if strings.HasPrefix(content, "**") && strings.HasSuffix(content, "**") {
		return &token{bold, content[2 : len(content)-2]}
	} else if strings.HasPrefix(content, "*") && strings.HasSuffix(content, "*") {
		return &token{italic, content[1 : len(content)-1]}
	} else {
		return &token{style, content}
	}
}

func inlineParseAndAppend(style int, content string) []*token {
	var line []*token

	InlineRE := fmt.Sprintf("%s|%s", BoldRE, ItalicRE)
	re := regexp.MustCompile(InlineRE)
	groups := re.FindAllStringIndex(content, -1)
	lastPos := 0

	for _, group := range groups {
		if group[0] > lastPos {
			line = append(line, newTokenDerived(content[lastPos:group[0]], style))
		}

		line = append(line, newTokenDerived(content[group[0]:group[1]], style))
		lastPos = group[1]
	}

	if lastPos < len(content) {
		line = append(line, newTokenDerived(content[lastPos:], style))
	}

	return line
}

// newParser returns a pointer to an object of mdParser.
// It parses the passed string into markdown tokens.
func newParser(input string) *mdParser {
	p := &mdParser{}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// headings
		var currTokens []*token

		r := regexp.MustCompile(HeadingRE)

		if r.MatchString(line) {
			level := len(r.FindStringSubmatch(line)[1])
			content := strings.TrimSpace(r.FindStringSubmatch(line)[2])
			/*
			* heading 1-6 has value 1-6 according to the enum declaration above.
			* So, we can directly use the length of `#` to determine the heading level
			 */
			currTokens = append(currTokens, &token{level, content})
			p.lines = append(p.lines, currTokens)

			continue
		}
		// paragraph
		currTokens = append(currTokens, inlineParseAndAppend(para, line)...)
		p.lines = append(p.lines, currTokens)
	}

	return p
}
