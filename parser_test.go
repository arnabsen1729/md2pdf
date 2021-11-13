package main // nolint

/*
	Naming the package for test can be either foo_test or just foo depending
	on if it is a black box testing or white box testing.

	More about it here: https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserOnSingleLineText(t *testing.T) {
	t.Parallel()

	// testing simple string as paragraph.
	p := newParser("hello")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: para, content: "hello", altContent: ""}, p.lines[0][0])

	// heading 1.
	p = newParser("# hello this is great")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: heading1, content: "hello this is great", altContent: ""}, p.lines[0][0])

	// heading 1 with extra spaces in front should not be treated as heading.
	p = newParser("   #   hello this is great    ")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: para, content: "#   hello this is great", altContent: ""}, p.lines[0][0])

	// heading 2.
	p = newParser("## hello")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: heading2, content: "hello", altContent: ""}, p.lines[0][0])
}

func TestParserOnMultiLineText(t *testing.T) {
	t.Parallel()

	// testing simple string as paragraph.
	p := newParser("hello\nthis is another line")
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, &token{style: para, content: "hello", altContent: ""}, p.lines[0][0])
	assert.Equal(t, &token{style: para, content: "this is another line", altContent: ""}, p.lines[1][0])

	// heading 1 with string
	p = newParser("# hello\nthis is great")
	// expP := [[&token{style: heading1, content: "hello"}], [&token{style: para, content: "this is great"}]]
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, []*token{{style: heading1, content: "hello", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{{style: para, content: "this is great", altContent: ""}}, p.lines[1])

	// heading 1 with extra spaces in front should be treated as para
	p = newParser("   #   hello\n     this \n is great    ")
	assert.Equal(t, 3, len(p.lines))
	assert.Equal(t, []*token{{style: para, content: "#   hello", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{{style: para, content: "this", altContent: ""}}, p.lines[1])
	assert.Equal(t, []*token{{style: para, content: "is great", altContent: ""}}, p.lines[2])

	// heading 2.
	p = newParser("#hello\n## hello2\n para ")
	assert.Equal(t, 3, len(p.lines))
	assert.Equal(t, []*token{{style: heading1, content: "hello", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{{style: heading2, content: "hello2", altContent: ""}}, p.lines[1])
	assert.Equal(t, []*token{{style: para, content: "para", altContent: ""}}, p.lines[2])
}

func TestParserOnInline(t *testing.T) {
	t.Parallel()

	// bold text
	p := newParser("**bold word**")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: bold, content: "bold word", altContent: ""}}, p.lines[0])

	// italic text
	p = newParser("*italic word*")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: italic, content: "italic word", altContent: ""}}, p.lines[0])

	// code text
	p = newParser("`code word`")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: code, content: "code word", altContent: ""}}, p.lines[0])

	// mixed
	p = newParser("normal **bold word** test with *italic word* and * is used with `code word`")
	assert.Equal(t, 1, len(p.lines))

	expTokens := []*token{
		{style: para, content: "normal", altContent: ""},
		{style: bold, content: "bold word", altContent: ""},
		{style: para, content: "test with", altContent: ""},
		{style: italic, content: "italic word", altContent: ""},
		{style: para, content: "and * is used with", altContent: ""},
		{style: code, content: "code word", altContent: ""},
	}
	assert.Equal(t, expTokens, p.lines[0])
}

func TestParserOnLink(t *testing.T) {
	t.Parallel()

	p := newParser("[link](https://www.google.com)")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: link, altContent: "https://www.google.com", content: "link"}}, p.lines[0])

	p = newParser("Normal Text with [link](https://www.google.com)")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{
		{style: para, content: "Normal Text with", altContent: ""},
		{style: link, altContent: "https://www.google.com", content: "link"},
	}, p.lines[0])

	p = newParser("New Line Text\n[link](https://www.google.com)")
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, []*token{{style: para, content: "New Line Text", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{
		{style: link, altContent: "https://www.google.com", content: "link"},
	}, p.lines[1])
}

func TestParserOnImage(t *testing.T) {
	t.Parallel()

	p := newParser("![link](https://cdn.recast.ai/newsletter/city-01.png)")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: image, altContent: "https://cdn.recast.ai/newsletter/city-01.png", content: "link"}}, p.lines[0])
}

func TestParseOnBlockQuotes(t *testing.T) {
	t.Parallel()

	p := newParser("> hello\n> this is another line")
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, []*token{{style: blockquote, content: "hello", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{{style: blockquote, content: "this is another line", altContent: ""}}, p.lines[1])

	p = newParser("# heading\n> this is bq\nthis is a line")
	assert.Equal(t, 3, len(p.lines))
	assert.Equal(t, []*token{{style: heading1, content: "heading", altContent: ""}}, p.lines[0])
	assert.Equal(t, []*token{{style: blockquote, content: "this is bq", altContent: ""}}, p.lines[1])
	assert.Equal(t, []*token{{style: para, content: "this is a line", altContent: ""}}, p.lines[2])
}
