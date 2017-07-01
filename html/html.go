package html

import (
	"errors"
	"fmt"

	"github.com/mokelab-go/markdown"
)

const (
	block_none = iota
	h1
	h2
	h3
	h4
	h5
	h6
	ul
	p
	code_block

	state_none = iota
	read_head_space

	read_bq_1
	read_bq_2
	read_bq_3
	code
	read_bq_end_1
	read_bq_end_2
	read_bq_end_3

	read_ul
	read_h1
	read_h2
	text
	text_link_start
	text_link_end
	text_link_url_start
	text_link_url_end
	text_image_1
	text_image_start
	text_image_end
	text_image_url_start
	text_image_url_end
	text_br
)

type impl struct {
}

func NewMarkdown() markdown.Markdown {
	return &impl{}
}

func (o *impl) Compile(src string) (string, error) {
	out := make([]byte, 0, len(src)*2)
	var linkText []byte
	var urlText []byte
	state := state_none
	block := block_none
	index := 0
	srcLen := len(src)
	for index < srcLen {
		char := src[index]
		switch state {
		case state_none:
			if char == '#' {
				state = read_h1
			} else if char == '`' {
				state = read_bq_1
			} else if char == '\n' {
				if block == ul {
					out = addBlockEnd(out, block)
					block = block_none
				}
				// ignore
			} else if char == ' ' {
				state = read_head_space
			} else {
				state = text
				block = p
				out = appendStr(out, "<p>")
				index--
			}
		case read_bq_1:
			if char == '`' {
				state = read_bq_2
			} else {
				state = text
				index--
			}
		case read_bq_2:
			if char == '`' {
				state = read_bq_3
			} else {
				state = text
				out = appendStr(out, "``")
				index--
			}
		case read_bq_3:
			if char == '\n' {
				state = code
				block = code_block
				out = appendStr(out, "<pre><code>")
			} else {
				// TODO : check lang
				state = text
				out = appendStr(out, "```")
				index--
			}
		case code:
			if char == '`' {
				state = read_bq_end_1
			} else {
				out = append(out, char)
			}
		case read_bq_end_1:
			if char == '`' {
				state = read_bq_end_2
			} else {
				out = appendStr(out, "`")
				index--
			}
		case read_bq_end_2:
			if char == '`' {
				state = read_bq_end_3
			} else {
				out = appendStr(out, "`1")
				index--
			}
		case read_bq_end_3:
			if char == '\n' {
				out = appendStr(out, "</code></pre>\n\n")
				block = block_none
				state = state_none
			} else {
				return "", errors.New(fmt.Sprintf("\\n expected byt %s at %d", char, index))
			}
		case read_head_space:
			if char == '*' || char == '-' {
				state = read_ul
			} else {
				state = text
				block = p
				out = appendStr(out, "<p>")
				out = append(out, char)
			}
		case read_ul:
			if char == ' ' {
				state = text
				if block == ul {
					out = appendStr(out, " <li>")
				} else {
					out = appendStr(out, "<ul>\n <li>")
					block = ul
				}
			} else {
				return "", errors.New(fmt.Sprintf("unexpected token at %d", index))
			}
		case text:
			if char == '\n' {
				if block == ul {
					state = state_none
					out = appendStr(out, "</li>\n")
				} else {
					state = text_br
				}
			} else if char == '[' {
				state = text_link_start
				linkText = make([]byte, 0)
			} else if char == '!' {
				state = text_image_1
				linkText = make([]byte, 0)
			} else if char == '<' {
				out = appendStr(out, "&lt;")
			} else if char == '>' {
				out = appendStr(out, "&gt;")
			} else {
				out = append(out, char)
			}
		case text_link_start:
			if char == ']' {
				state = text_link_end
			} else {
				linkText = append(linkText, char)
			}
		case text_link_end:
			if char == '(' {
				state = text_link_url_start
				urlText = make([]byte, 0)
			} else {
				return "", errors.New(fmt.Sprintf("expected is ( but %s at %d", char, index))
			}
		case text_link_url_start:
			if char == ')' {
				state = text
				out = appendStr(out, fmt.Sprintf("<a href=\"%s\">%s</a>", string(urlText), string(linkText)))
			} else {
				urlText = append(urlText, char)
			}
		case text_image_1:
			if char == '[' {
				state = text_image_start
			} else {
				out = appendStr(out, "!")
				index--
				state = text
			}
		case text_image_start:
			if char == ']' {
				state = text_image_end
			} else {
				linkText = append(linkText, char)
			}
		case text_image_end:
			if char == '(' {
				state = text_image_url_start
				urlText = make([]byte, 0)
			} else {
				return "", errors.New(fmt.Sprintf("expected is ( but %s at %d", char, index))
			}
		case text_image_url_start:
			if char == ')' {
				state = text
				out = appendStr(out, fmt.Sprintf("<img src=\"%s\" title=\"%s\"/>", string(urlText), string(linkText)))
			} else {
				urlText = append(urlText, char)
			}
		case text_br:
			if char == '\n' {
				state = state_none
				out = addBlockEnd(out, block)
				block = block_none
			} else {
				state = text
				out = append(out, ' ')
				out = append(out, char)
			}
		case read_h1: // #
			if char == ' ' {
				block = h1
				state = text
				out = appendStr(out, "<h1>")
			} else if char == '#' {
				state = read_h2
			}
		case read_h2: // ##
			if char == ' ' {
				block = h2
				state = text
				out = appendStr(out, "<h2>")
			}
		}
		index++
	}
	out = addBlockEnd(out, block)

	return string(out), nil
}

func addBlockEnd(out []byte, block int) []byte {
	switch block {
	case h1:
		out = appendStr(out, "</h1>\n\n")
	case h2:
		out = appendStr(out, "</h2>\n\n")
	case ul:
		out = appendStr(out, "</ul>\n\n")
	case p:
		out = appendStr(out, "</p>\n\n")
	}
	return out
}

func appendStr(out []byte, text string) []byte {
	return append(out, text...)
}
