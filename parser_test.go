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

func TestParserOnText(t *testing.T) {
	// testing simple string as paragraph.
	p := newParser("hello")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: para, content: "hello"})

	// heading 1.
	p = newParser("# hello this is great")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: heading1, content: "hello this is great"})

	// heading 1 with extra spaces.
	p = newParser("   #   hello this is great    ")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: heading1, content: "hello this is great"})

	// heading 2.
	p = newParser("## hello")
	assert.NotEmpty(t, p.tokens)
	assert.Equal(t, p.tokens[0], &token{style: heading2, content: "hello"})
}
