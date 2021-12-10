package parser // nolint

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
		p := NewParser("hello")
		assert.NotEmpty(t, p.Lines)
		assert.Equal(t, &Token{Style: Para, Content: "hello", AltContent: ""}, p.Lines[0][0])
	})

	t.Run("heading 1", func(t *testing.T) {
		t.Parallel()
		p := NewParser("# hello this is great")
		assert.NotEmpty(t, p.Lines)
		assert.Equal(t, &Token{Style: Heading1, Content: "hello this is great", AltContent: ""}, p.Lines[0][0])
	})

	t.Run("heading 1 with extra spaces in front should not be treated as heading", func(t *testing.T) {
		t.Parallel()
		p := NewParser("   #   hello this is great    ")
		assert.NotEmpty(t, p.Lines)
		assert.Equal(t, &Token{Style: Para, Content: "#   hello this is great", AltContent: ""}, p.Lines[0][0])
	})

	t.Run("heading 2", func(t *testing.T) {
		t.Parallel()
		p := NewParser("## hello")
		assert.NotEmpty(t, p.Lines)
		assert.Equal(t, &Token{Style: Heading2, Content: "hello", AltContent: ""}, p.Lines[0][0])
	})
}

func TestParserOnMultiLineText(t *testing.T) {
	t.Parallel()

	t.Run("testing simple string as paragraph", func(t *testing.T) {
		t.Parallel()
		p := NewParser("hello\nthis is another line")
		assert.Equal(t, 2, len(p.Lines))
		assert.Equal(t, &Token{Style: Para, Content: "hello", AltContent: ""}, p.Lines[0][0])
		assert.Equal(t, &Token{Style: Para, Content: "this is another line", AltContent: ""}, p.Lines[1][0])
	})

	t.Run("heading 1 with string", func(t *testing.T) {
		t.Parallel()
		p := NewParser("# hello\nthis is great")
		// expP := [[&token{style: heading1, content: "hello"}], [&token{style: para, content: "this is great"}]]
		assert.Equal(t, 2, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Heading1, Content: "hello", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{{Style: Para, Content: "this is great", AltContent: ""}}, p.Lines[1])
	})

	t.Run("heading 1 with extra spaces in front should be treated as para", func(t *testing.T) {
		t.Parallel()
		p := NewParser("   #   hello\n     this \n is great    ")
		assert.Equal(t, 3, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Para, Content: "#   hello", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{{Style: Para, Content: "this", AltContent: ""}}, p.Lines[1])
		assert.Equal(t, []*Token{{Style: Para, Content: "is great", AltContent: ""}}, p.Lines[2])
	})

	t.Run("heading 1, heading 2 and a paragraph", func(t *testing.T) {
		t.Parallel()
		p := NewParser("#hello\n## hello2\n para ")
		assert.Equal(t, 3, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Heading1, Content: "hello", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{{Style: Heading2, Content: "hello2", AltContent: ""}}, p.Lines[1])
		assert.Equal(t, []*Token{{Style: Para, Content: "para", AltContent: ""}}, p.Lines[2])
	})
}

func TestParserOnInline(t *testing.T) {
	t.Parallel()

	t.Run("bold text", func(t *testing.T) {
		t.Parallel()
		p := NewParser("**bold word**")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Bold, Content: "bold word", AltContent: ""}}, p.Lines[0])
	})

	t.Run("italic text", func(t *testing.T) {
		t.Parallel()
		p := NewParser("*italic word*")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Italic, Content: "italic word", AltContent: ""}}, p.Lines[0])
	})

	t.Run("code text", func(t *testing.T) {
		t.Parallel()
		p := NewParser("`code word`")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Code, Content: "code word", AltContent: ""}}, p.Lines[0])
	})

	t.Run("mixed", func(t *testing.T) {
		t.Parallel()
		p := NewParser("normal **bold word** test with *italic word* and * is used with `code word`")
		assert.Equal(t, 1, len(p.Lines))

		expTokens := []*Token{
			{Style: Para, Content: "normal", AltContent: ""},
			{Style: Bold, Content: "bold word", AltContent: ""},
			{Style: Para, Content: "test with", AltContent: ""},
			{Style: Italic, Content: "italic word", AltContent: ""},
			{Style: Para, Content: "and * is used with", AltContent: ""},
			{Style: Code, Content: "code word", AltContent: ""},
		}
		assert.Equal(t, expTokens, p.Lines[0])
	})
}

func TestParserOnLink(t *testing.T) {
	t.Parallel()

	t.Run("one word link", func(t *testing.T) {
		t.Parallel()
		p := NewParser("[link](https://www.google.com)")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Link, AltContent: "https://www.google.com", Content: "link"}}, p.Lines[0])
	})

	t.Run("inline link with text", func(t *testing.T) {
		t.Parallel()
		p := NewParser("Normal Text with [link](https://www.google.com)")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{
			{Style: Para, Content: "Normal Text with", AltContent: ""},
			{Style: Link, AltContent: "https://www.google.com", Content: "link"},
		}, p.Lines[0])
	})

	t.Run("text with link in newline", func(t *testing.T) {
		t.Parallel()
		p := NewParser("New Line Text\n[link](https://www.google.com)")
		assert.Equal(t, 2, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Para, Content: "New Line Text", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{
			{Style: Link, AltContent: "https://www.google.com", Content: "link"},
		}, p.Lines[1])
	})
}

func TestParserOnImage(t *testing.T) {
	t.Parallel()

	t.Run("image", func(t *testing.T) {
		t.Parallel()
		p := NewParser("![link](https://cdn.recast.ai/newsletter/city-01.png)")
		assert.Equal(t, 1, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Image, AltContent: "https://cdn.recast.ai/newsletter/city-01.png", Content: "link"}}, p.Lines[0])
	})
}

func TestParseOnBlockQuotes(t *testing.T) {
	t.Parallel()

	t.Run("consecutive line blockquotes", func(t *testing.T) {
		t.Parallel()
		p := NewParser("> hello\n> this is another line")
		assert.Equal(t, 2, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Blockquote, Content: "hello", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{{Style: Blockquote, Content: "this is another line", AltContent: ""}}, p.Lines[1])
	})

	t.Run("blockquotes with heading and text", func(t *testing.T) {
		t.Parallel()
		p := NewParser("# heading\n> this is bq\nthis is a line")
		assert.Equal(t, 3, len(p.Lines))
		assert.Equal(t, []*Token{{Style: Heading1, Content: "heading", AltContent: ""}}, p.Lines[0])
		assert.Equal(t, []*Token{{Style: Blockquote, Content: "this is bq", AltContent: ""}}, p.Lines[1])
		assert.Equal(t, []*Token{{Style: Para, Content: "this is a line", AltContent: ""}}, p.Lines[2])
	})
}
