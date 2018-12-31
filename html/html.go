package html

import (
	"fmt"

	"github.com/mokelab-go/markdown"
	"github.com/mokelab-go/markdown/ast"
)

type impl struct {
}

func NewMarkdown() markdown.Markdown {
	return &impl{}
}

func (o *impl) Compile(src string) (string, error) {
	tree, err := ast.Parse(src)
	if err != nil {
		return "", err
	}
	out := make([]byte, 0, len(src)*2)
	out = printBlock(out, tree)
	return string(out), nil
}

func appendStr(out []byte, text string) []byte {
	return append(out, text...)
}

func printBlock(out []byte, block *ast.Block) []byte {
	switch block.Type {
	case ast.TypeRoot:
		return printChildren(out, block)
	case ast.TypeH1:
		out = appendStr(out, "<h1>")
		out = printChildren(out, block)
		out = appendStr(out, "</h1>\n\n")
	case ast.TypeH2:
		out = appendStr(out, "<h2>")
		out = printChildren(out, block)
		out = appendStr(out, "</h2>\n\n")
	case ast.TypeP:
		out = appendStr(out, "<p>")
		out = printChildren(out, block)
		out = appendStr(out, "</p>\n\n")
	case ast.TypePreCode:
		out = appendStr(out, "<pre><code>")
		out = printChildren(out, block)
		out = appendStr(out, "</code></pre>\n\n")
	case ast.TypeUL:
		out = appendStr(out, "<ul>\n")
		out = printChildren(out, block)
		out = appendStr(out, "</ul>\n\n")
	case ast.TypeLI:
		out = appendStr(out, " <li>")
		out = printChildren(out, block)
		out = appendStr(out, " </li>\n")
	case ast.TypeAnchor:
		out = appendStr(out, fmt.Sprintf("<a href=\"%s\">%s</a>", block.URL, block.Value))
	case ast.TypeImage:
		out = appendStr(out, fmt.Sprintf("<img src=\"%s\" title=\"%s\"/>", block.URL, block.Value))
	case ast.TypeText:
		if len(block.Value) == 0 {
			out = printChildren(out, block)
		} else {
			out = appendStr(out, block.Value)
		}
	case ast.TypeCode:
		out = appendStr(out, "<code>")
		out = appendStr(out, block.Value)
		out = appendStr(out, "</code>\n\n")
	}
	return out
}

func printChildren(out []byte, block *ast.Block) []byte {
	for _, e := range block.Children {
		out = printBlock(out, e)
	}
	return out
}
