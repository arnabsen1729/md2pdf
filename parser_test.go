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
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: para, content: "hello"})

	// heading 1.
	p = newParser("# hello this is great")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: heading1, content: "hello this is great"})

	// heading 1 with extra spaces in front should not be treated as heading.
	p = newParser("   #   hello this is great    ")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, &token{style: para, content: "   #   hello this is great    "}, p.tokens[0])

	// heading 2.
	p = newParser("## hello")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, &token{style: heading2, content: "hello"}, p.tokens[0])
}

func TestParserOnMultiLineText(t *testing.T) {
	// testing simple string as paragraph.
	p := newParser("hello\nthis is another line")
	assert.Equal(t, len(p.tokens), 2)
	assert.Equal(t, &token{style: para, content: "hello"}, p.tokens[0])
	assert.Equal(t, &token{style: para, content: "this is another line"}, p.tokens[1])

	// heading 1 with string
	p = newParser("# hello\nthis is great")

	assert.Equal(t, len(p.tokens), 2)
	assert.Equal(t, &token{style: heading1, content: "hello"}, p.tokens[0])
	assert.Equal(t, &token{style: para, content: "this is great"}, p.tokens[1])

	// heading 1 with extra spaces in front should be treated as para
	p = newParser("   #   hello\n     this \n is great    ")
	assert.Equal(t, len(p.tokens), 3)
	assert.Equal(t, &token{style: para, content: "   #   hello"}, p.tokens[0])
	assert.Equal(t, &token{style: para, content: "     this "}, p.tokens[1])
	assert.Equal(t, &token{style: para, content: " is great    "}, p.tokens[2])

	// heading 2.
	p = newParser("#hello\n## hello2\n para ")
	assert.Equal(t, len(p.tokens), 3)
	assert.Equal(t, &token{style: heading1, content: "hello"}, p.tokens[0])
	assert.Equal(t, &token{style: heading2, content: "hello2"}, p.tokens[1])
	assert.Equal(t, &token{style: para, content: " para "}, p.tokens[2])

}
