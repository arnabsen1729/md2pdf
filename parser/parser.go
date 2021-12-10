// Package parser will parse the markdown files into tokens. Each token has it's
// own specific style. We pass the text to the NewParser() function which
// returns an MdParser object pointer.
package parser

import (
	"regexp"
	"strings"
)

/*
	There is no default structure in Golang called enums (enumerators).
	We can simulate it using constant variable declaration.
	Link: https://www.sohamkamani.com/golang/enums/
*/

// Markdown styles.
const (
	Para = iota
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

// Regex to match the respective tags
// Regex source: https://github.com/Python-Markdown/markdown/blob/master/markdown/blockprocessors.py#L448
const (
	HeadingRE    = `(?:^|\n)(?P<level>#{1,6})(?P<header>(?:\\.|[^\\])*?)#*(?:\n|$)`
	BoldRE       = `(\*{2}.+?\*{2})`
	ItalicRE     = `(\*{1}.+?\*{1})`
	CodeRE       = `(\` + "`" + `{1}.+?\` + "`" + `{1})` // Golang doesn't support ` (backtick) in raw strings.
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

// Token is the basic unit of a markdown file. For links, AltContent will store
// the URL and content will store the text to display. Similarly, for images
// content will store the URL of the image.
type Token struct {
	Style      int
	Content    string
	AltContent string // only in links and images, for rest it will ""(empty string).
}

// MdParser will parse the text passed into markdown tokens. Each line in a
// markdown file can be represented as a list of tokens (inline styling). So the
// entire markdown file is a list of such lines and hence a list of list of
// token.
type MdParser struct {
	Lines [][]*Token
}

func newTokenDerived(content string, style int) *Token {
	content = strings.TrimSpace(content)

	switch {
	case strings.HasPrefix(content, "**") && strings.HasSuffix(content, "**"): // bold
		return &Token{Style: Bold, Content: content[2 : len(content)-2], AltContent: ""}
	case strings.HasPrefix(content, "*") && strings.HasSuffix(content, "*"): // italic
		return &Token{Style: Italic, Content: content[1 : len(content)-1], AltContent: ""}
	case strings.HasPrefix(content, "`") && strings.HasSuffix(content, "`"): // code (inline)
		return &Token{Style: Code, Content: content[1 : len(content)-1], AltContent: ""}
	case strings.HasPrefix(content, "[") && strings.HasSuffix(content, ")"): // link
		closingBracketPos := strings.Index(content, "]")
		linkContent := content[1:closingBracketPos]
		linkURL := content[closingBracketPos+2 : len(content)-1]

		return &Token{Style: Link, Content: linkContent, AltContent: linkURL}
	case strings.HasPrefix(content, "![") && strings.HasSuffix(content, ")"): // image
		closingBracketPos := strings.Index(content, "]")
		imageContent := content[2:closingBracketPos]
		imageURL := content[closingBracketPos+2 : len(content)-1]

		return &Token{Style: Image, Content: imageContent, AltContent: imageURL}
	default:
		return &Token{Style: style, Content: content, AltContent: ""}
	}
}

func inlineParseAndAppend(style int, content string) []*Token {
	line := make([]*Token, 0)

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

// NewParser returns a pointer to an object of MdParser.
// It parses the passed string into markdown tokens.
func NewParser(input string) *MdParser {
	p := &MdParser{Lines: nil}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// headings
		var currTokens []*Token

		switch {
		case headingCompiledRE.MatchString(line): // heading
			level := len(headingCompiledRE.FindStringSubmatch(line)[1])
			content := strings.TrimSpace(headingCompiledRE.FindStringSubmatch(line)[2])
			/*
			* heading 1-6 has value 1-6 according to the enum declaration above.
			* So, we can directly use the length of `#` to determine the heading level
			 */
			currTokens = append(currTokens, &Token{Style: level, Content: content, AltContent: ""})
		case blockquoteCompiledRE.MatchString(line): // blockquote
			content := strings.TrimSpace(blockquoteCompiledRE.FindStringSubmatch(line)[2])

			currTokens = append(currTokens, &Token{Style: Blockquote, Content: content, AltContent: ""})
		default: // paragraph
			currTokens = append(currTokens, inlineParseAndAppend(Para, line)...)
		}

		p.Lines = append(p.Lines, currTokens)
	}

	return p
}
