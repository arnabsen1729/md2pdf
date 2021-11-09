# Markdown to PDF

![hero](/.github/assets/hero.png)

Will take a markdown file as input and then create a PDF file with the markdown formatting.

## Usage

```
Usage of md2pdf:
  -file string
    	Name of the markdown file to read
  -output string
    	Name of the PDF file to be exported  (default: <input-file-name>.pdf)
```

Example:

```
$ md2pdf -file=MyFile.md -output=MyFile.pdf
```

## Parser

A Parser will parse the input markdown file as tokens.

- [X] Headings (L1 - L6)
- [X] Paragraph
- [ ] Blockquotes
- [ ] CodeBlock
- [ ] Lists (Ordered and Unordered)
- [ ] Links
- [ ] Horizontal Rules
- [ ] Tables

Will also consider some inline blocks like:

- [X] Bold
- [X] Italic
- [X] Code
- [ ] Link
- [ ] Images

[Markdown Guidelines](https://www.markdownguide.org/basic-syntax/)
