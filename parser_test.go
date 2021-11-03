package main

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
	// testing simple string as paragraph.
	p := newParser("hello")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: para, content: "hello"}, p.lines[0][0])

	// heading 1.
	p = newParser("# hello this is great")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: heading1, content: "hello this is great"}, p.lines[0][0])

	// heading 1 with extra spaces in front should not be treated as heading.
	p = newParser("   #   hello this is great    ")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: para, content: "#   hello this is great"}, p.lines[0][0])

	// heading 2.
	p = newParser("## hello")
	assert.NotEmpty(t, p.lines)
	assert.Equal(t, &token{style: heading2, content: "hello"}, p.lines[0][0])
}

func TestParserOnMultiLineText(t *testing.T) {
	// testing simple string as paragraph.
	p := newParser("hello\nthis is another line")
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, &token{style: para, content: "hello"}, p.lines[0][0])
	assert.Equal(t, &token{style: para, content: "this is another line"}, p.lines[1][0])

	// heading 1 with string
	p = newParser("# hello\nthis is great")
	// expP := [[&token{style: heading1, content: "hello"}], [&token{style: para, content: "this is great"}]]
	assert.Equal(t, 2, len(p.lines))
	assert.Equal(t, []*token{{style: heading1, content: "hello"}}, p.lines[0])
	assert.Equal(t, []*token{{style: para, content: "this is great"}}, p.lines[1])

	// heading 1 with extra spaces in front should be treated as para
	p = newParser("   #   hello\n     this \n is great    ")
	assert.Equal(t, 3, len(p.lines))
	assert.Equal(t, []*token{{style: para, content: "#   hello"}}, p.lines[0])
	assert.Equal(t, []*token{{style: para, content: "this"}}, p.lines[1])
	assert.Equal(t, []*token{{style: para, content: "is great"}}, p.lines[2])

	// heading 2.
	p = newParser("#hello\n## hello2\n para ")
	assert.Equal(t, 3, len(p.lines))
	assert.Equal(t, []*token{{style: heading1, content: "hello"}}, p.lines[0])
	assert.Equal(t, []*token{{style: heading2, content: "hello2"}}, p.lines[1])
	assert.Equal(t, []*token{{style: para, content: "para"}}, p.lines[2])
}

func TestParserOnInline(t *testing.T) {
	// bold text
	p := newParser("**bold**")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: bold, content: "bold"}}, p.lines[0])

	// italic text
	p = newParser("*italic*")
	assert.Equal(t, 1, len(p.lines))
	assert.Equal(t, []*token{{style: italic, content: "italic"}}, p.lines[0])

	// mixed
	p = newParser("normal **bold** test with *italic* and * is used")
	assert.Equal(t, 1, len(p.lines))
	expTokens := []*token{
		{style: para, content: "normal"},
		{style: bold, content: "bold"},
		{style: para, content: "test with"},
		{style: italic, content: "italic"},
		{style: para, content: "and * is used"},
	}
	assert.Equal(t, expTokens, p.lines[0])
}
