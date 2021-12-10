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

	t.Run("simple string as paragraph", func(t *testing.T) {
		t.Parallel()
		p := newParser("hello")
		assert.NotEmpty(t, p.lines)
		assert.Equal(t, &token{style: Para, content: "hello", altContent: ""}, p.lines[0][0])
	})

	t.Run("heading 1", func(t *testing.T) {
		t.Parallel()
		p := newParser("# hello this is great")
		assert.NotEmpty(t, p.lines)
		assert.Equal(t, &token{style: Heading1, content: "hello this is great", altContent: ""}, p.lines[0][0])
	})

	t.Run("heading 1 with extra spaces in front should not be treated as heading", func(t *testing.T) {
		t.Parallel()
		p := newParser("   #   hello this is great    ")
		assert.NotEmpty(t, p.lines)
		assert.Equal(t, &token{style: Para, content: "#   hello this is great", altContent: ""}, p.lines[0][0])
	})

	t.Run("heading 2", func(t *testing.T) {
		t.Parallel()
		p := newParser("## hello")
		assert.NotEmpty(t, p.lines)
		assert.Equal(t, &token{style: Heading2, content: "hello", altContent: ""}, p.lines[0][0])
	})
}

func TestParserOnMultiLineText(t *testing.T) {
	t.Parallel()

	t.Run("testing simple string as paragraph", func(t *testing.T) {
		t.Parallel()
		p := newParser("hello\nthis is another line")
		assert.Equal(t, 2, len(p.lines))
		assert.Equal(t, &token{style: Para, content: "hello", altContent: ""}, p.lines[0][0])
		assert.Equal(t, &token{style: Para, content: "this is another line", altContent: ""}, p.lines[1][0])
	})

	t.Run("heading 1 with string", func(t *testing.T) {
		t.Parallel()
		p := newParser("# hello\nthis is great")
		// expP := [[&token{style: heading1, content: "hello"}], [&token{style: para, content: "this is great"}]]
		assert.Equal(t, 2, len(p.lines))
		assert.Equal(t, []*token{{style: Heading1, content: "hello", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{{style: Para, content: "this is great", altContent: ""}}, p.lines[1])
	})

	t.Run("heading 1 with extra spaces in front should be treated as para", func(t *testing.T) {
		t.Parallel()
		p := newParser("   #   hello\n     this \n is great    ")
		assert.Equal(t, 3, len(p.lines))
		assert.Equal(t, []*token{{style: Para, content: "#   hello", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{{style: Para, content: "this", altContent: ""}}, p.lines[1])
		assert.Equal(t, []*token{{style: Para, content: "is great", altContent: ""}}, p.lines[2])
	})

	t.Run("heading 1, heading 2 and a paragraph", func(t *testing.T) {
		t.Parallel()
		p := newParser("#hello\n## hello2\n para ")
		assert.Equal(t, 3, len(p.lines))
		assert.Equal(t, []*token{{style: Heading1, content: "hello", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{{style: Heading2, content: "hello2", altContent: ""}}, p.lines[1])
		assert.Equal(t, []*token{{style: Para, content: "para", altContent: ""}}, p.lines[2])
	})
}

func TestParserOnInline(t *testing.T) {
	t.Parallel()

	t.Run("bold text", func(t *testing.T) {
		t.Parallel()
		p := newParser("**bold word**")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{{style: Bold, content: "bold word", altContent: ""}}, p.lines[0])
	})

	t.Run("italic text", func(t *testing.T) {
		t.Parallel()
		p := newParser("*italic word*")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{{style: Italic, content: "italic word", altContent: ""}}, p.lines[0])
	})

	t.Run("code text", func(t *testing.T) {
		t.Parallel()
		p := newParser("`code word`")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{{style: Code, content: "code word", altContent: ""}}, p.lines[0])
	})

	t.Run("mixed", func(t *testing.T) {
		t.Parallel()
		p := newParser("normal **bold word** test with *italic word* and * is used with `code word`")
		assert.Equal(t, 1, len(p.lines))

		expTokens := []*token{
			{style: Para, content: "normal", altContent: ""},
			{style: Bold, content: "bold word", altContent: ""},
			{style: Para, content: "test with", altContent: ""},
			{style: Italic, content: "italic word", altContent: ""},
			{style: Para, content: "and * is used with", altContent: ""},
			{style: Code, content: "code word", altContent: ""},
		}
		assert.Equal(t, expTokens, p.lines[0])
	})
}

func TestParserOnLink(t *testing.T) {
	t.Parallel()

	t.Run("one word link", func(t *testing.T) {
		t.Parallel()
		p := newParser("[link](https://www.google.com)")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{{style: Link, altContent: "https://www.google.com", content: "link"}}, p.lines[0])
	})

	t.Run("inline link with text", func(t *testing.T) {
		t.Parallel()
		p := newParser("Normal Text with [link](https://www.google.com)")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{
			{style: Para, content: "Normal Text with", altContent: ""},
			{style: Link, altContent: "https://www.google.com", content: "link"},
		}, p.lines[0])
	})

	t.Run("text with link in newline", func(t *testing.T) {
		t.Parallel()
		p := newParser("New Line Text\n[link](https://www.google.com)")
		assert.Equal(t, 2, len(p.lines))
		assert.Equal(t, []*token{{style: Para, content: "New Line Text", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{
			{style: Link, altContent: "https://www.google.com", content: "link"},
		}, p.lines[1])
	})
}

func TestParserOnImage(t *testing.T) {
	t.Parallel()

	t.Run("image", func(t *testing.T) {
		t.Parallel()
		p := newParser("![link](https://cdn.recast.ai/newsletter/city-01.png)")
		assert.Equal(t, 1, len(p.lines))
		assert.Equal(t, []*token{{style: Image, altContent: "https://cdn.recast.ai/newsletter/city-01.png", content: "link"}}, p.lines[0])
	})
}

func TestParseOnBlockQuotes(t *testing.T) {
	t.Parallel()

	t.Run("consecutive line blockquotes", func(t *testing.T) {
		t.Parallel()
		p := newParser("> hello\n> this is another line")
		assert.Equal(t, 2, len(p.lines))
		assert.Equal(t, []*token{{style: Blockquote, content: "hello", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{{style: Blockquote, content: "this is another line", altContent: ""}}, p.lines[1])
	})

	t.Run("blockquotes with heading and text", func(t *testing.T) {
		t.Parallel()
		p := newParser("# heading\n> this is bq\nthis is a line")
		assert.Equal(t, 3, len(p.lines))
		assert.Equal(t, []*token{{style: Heading1, content: "heading", altContent: ""}}, p.lines[0])
		assert.Equal(t, []*token{{style: Blockquote, content: "this is bq", altContent: ""}}, p.lines[1])
		assert.Equal(t, []*token{{style: Para, content: "this is a line", altContent: ""}}, p.lines[2])
	})
}
