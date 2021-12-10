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
	Para int = iota
	Heading1
	Heading2
	Heading3
	Heading4
	Heading5
	Heading6
	Bold
	Italic
	Code
	Link
	Image
	Blockquote
)

// Regex source: https://github.com/Python-Markdown/markdown/blob/master/markdown/blockprocessors.py#L448
const (
	HeadingRE = `(?:^|\n)(?P<level>#{1,6})(?P<header>(?:\\.|[^\\])*?)#*(?:\n|$)`
	BoldRE    = `(\*{2}.+?\*{2})`
	ItalicRE  = `(\*{1}.+?\*{1})`
	CodeRE    = `(\` + "`" + `{1}.+?\` + "`" + `{1})`
	// Golang doesn't seem to support ` (backtick) in raw strings.
	// Ref: https://github.com/golang/go/issues/18221#issuecomment-265314494
	LinkRE       = `\[.+?\]\(.+?\)`
	ImageRE      = `\!\[.+?\]\(.+?\)`
	BlockquoteRE = `(^|\n)[ ]{0,3}>[ ]?(.*)`
)

// * For testing Regex use regex101.com .

// OR-ing the regexes together, to catch all the inline styles.
var (
	re                   = regexp.MustCompile(BoldRE + "|" + ItalicRE + "|" + CodeRE + "|" + LinkRE + "|" + ImageRE)
	headingCompiledRE    = regexp.MustCompile(HeadingRE)
	blockquoteCompiledRE = regexp.MustCompile(BlockquoteRE)
)

/*
	For links, altContent will store the URL and content will store the text to display.
	Similarly, for images content will store the URL of the image.
*/

type token struct {
	style      int
	content    string
	altContent string // only in links and images, for rest it will ""(empty string).
}

// mdParser will parse the text passed into markdown tokens.
type mdParser struct {
	lines [][]*token
}

func newTokenDerived(content string, style int) *token {
	content = strings.TrimSpace(content)

	switch {
	case strings.HasPrefix(content, "**") && strings.HasSuffix(content, "**"): // bold
		return &token{style: Bold, content: content[2 : len(content)-2], altContent: ""}
	case strings.HasPrefix(content, "*") && strings.HasSuffix(content, "*"): // italic
		return &token{style: Italic, content: content[1 : len(content)-1], altContent: ""}
	case strings.HasPrefix(content, "`") && strings.HasSuffix(content, "`"): // code (inline)
		return &token{style: Code, content: content[1 : len(content)-1], altContent: ""}
	case strings.HasPrefix(content, "[") && strings.HasSuffix(content, ")"): // link
		closingBracketPos := strings.Index(content, "]")
		linkContent := content[1:closingBracketPos]
		linkURL := content[closingBracketPos+2 : len(content)-1]

		return &token{style: Link, content: linkContent, altContent: linkURL}
	case strings.HasPrefix(content, "![") && strings.HasSuffix(content, ")"): // image
		closingBracketPos := strings.Index(content, "]")
		imageContent := content[2:closingBracketPos]
		imageURL := content[closingBracketPos+2 : len(content)-1]

		return &token{style: Image, content: imageContent, altContent: imageURL}
	default:
		return &token{style: style, content: content, altContent: ""}
	}
}

func inlineParseAndAppend(style int, content string) []*token {
	line := make([]*token, 0)

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
	p := &mdParser{lines: nil}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// headings
		var currTokens []*token

		switch {
		case headingCompiledRE.MatchString(line): // heading
			level := len(headingCompiledRE.FindStringSubmatch(line)[1])
			content := strings.TrimSpace(headingCompiledRE.FindStringSubmatch(line)[2])
			/*
			* heading 1-6 has value 1-6 according to the enum declaration above.
			* So, we can directly use the length of `#` to determine the heading level
			 */
			currTokens = append(currTokens, &token{style: level, content: content, altContent: ""})
		case blockquoteCompiledRE.MatchString(line): // blockquote
			content := strings.TrimSpace(blockquoteCompiledRE.FindStringSubmatch(line)[2])

			currTokens = append(currTokens, &token{style: Blockquote, content: content, altContent: ""})
		default: // paragraph
			currTokens = append(currTokens, inlineParseAndAppend(Para, line)...)
		}

		p.lines = append(p.lines, currTokens)
	}

	return p
}
