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

	state_none = iota
	read_head_space
	read_ul
	read_h1
	read_h2
	text
	text_br
)

type impl struct {
}

func NewMarkdown() markdown.Markdown {
	return &impl{}
}

func (o *impl) Compile(src string) (string, error) {
	out := make([]byte, 0, len(src)*2)
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
				out = append(out, char)
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
			} else if char == '<' {
				out = appendStr(out, "&lt;")
			} else if char == '>' {
				out = appendStr(out, "&gt;")
			} else {
				out = append(out, char)
			}
		case text_br:
			if char == '\n' {
				state = state_none
				out = addBlockEnd(out, block)
				block = 0
			} else {
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
