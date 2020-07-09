package ast

import (
	"testing"
)

const src_1 = `
# Hello

World![image](./a.webp width=100 height=200)
`

const src_2 = `
# Hello

## World

`

const src_3 = `
Hey

 * a
 * b
 * [c](https://mokelab.com)

`

const src4 = "Hey\n\n" +
	"```\n" +
	"package main\n\n" +
	"func main() {\n" +
	"}\n" +
	"```\n\n" +
	"text node"

const src_5 = "Hey\n\n" +
	"```java\n" +
	"package main\n\n" +
	"func main() {\n" +
	"}\n" +
	"```\n"

const src_6 = "Hey\n\n" +
	"```\n" +
	"package main\n\n" +
	"func main() {\n" +
	"   s := `\n" +
	"`\n" +
	"}\n" +
	"```\n"

const src7 = "# Program\n" +
	"\n" +
	"This is `code` block\n"

func Test_Stack(t *testing.T) {
	stack := &blockStack{values: make([]*Block, 0)}
	stack.Push(newBlock(TypeRoot))
	stack.Push(newBlock(TypeP))
	v := stack.Pop()
	if v.Type != TypeP {
		t.Errorf("Wrong node")
	}
	v = stack.Pop()
	if v.Type != TypeRoot {
		t.Errorf("Wrong node")
	}
}

func Test_1(t *testing.T) {
	out, err := Parse(src_1)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- h1
	//    |   |- text
	//    |- p
	//       |- text
	//       |- image
	//       |- text(empty)
	checkBlock(t, out, TypeRoot, 2)

	h1Block := out.Children[0]
	checkBlock(t, h1Block, TypeH1, 1)

	h1Text := h1Block.Children[0]
	checkTextBlock(t, h1Text, "Hello")

	pBlock := out.Children[1]
	checkBlock(t, pBlock, TypeP, 3)

	text := pBlock.Children[0]
	checkTextBlock(t, text, "World")
	imageBlock := pBlock.Children[1]
	checkImageBlock(t, imageBlock, "image", "./a.webp")
	if len(imageBlock.Attributes) != 2 {
		t.Errorf("Attributes must have 2 but %d", len(imageBlock.Attributes))
		return
	}
	v, ok := imageBlock.Attributes["width"]
	if !ok {
		t.Errorf("Attributes must have width")
		return
	}
	if v != "100" {
		t.Errorf("Width must be 100 but %s", v)
		return
	}
	v, ok = imageBlock.Attributes["height"]
	if !ok {
		t.Errorf("Attributes must have height")
		return
	}
	if v != "200" {
		t.Errorf("height must be 200 but %s", v)
		return
	}
	text = pBlock.Children[2]
	checkTextBlock(t, text, "")
}

func Test_2(t *testing.T) {
	out, err := Parse(src_2)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- h1
	//    |   |- text
	//    |- h2
	//       |- text
	checkBlock(t, out, TypeRoot, 2)

	h1Block := out.Children[0]
	checkBlock(t, h1Block, TypeH1, 1)

	h1Text := h1Block.Children[0]
	checkTextBlock(t, h1Text, "Hello")

	h2Block := out.Children[1]
	checkBlock(t, h2Block, TypeH2, 1)

	h2Text := h2Block.Children[0]
	checkTextBlock(t, h2Text, "World")
}

func Test_3(t *testing.T) {
	out, err := Parse(src_3)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- p
	//    |   |- text
	//    |- ul
	//       |- li
	//       |- li
	//       |- li
	//           |- text(empty)
	//           |- anchor
	//           |- text(empty)
	checkBlock(t, out, TypeRoot, 2)

	pBlock := out.Children[0]
	checkBlock(t, pBlock, TypeP, 1)

	pText := pBlock.Children[0]
	checkTextBlock(t, pText, "Hey")

	ulBlock := out.Children[1]
	checkBlock(t, ulBlock, TypeUL, 3)

	liBlock := ulBlock.Children[0]
	checkBlock(t, liBlock, TypeLI, 1)
	liText := liBlock.Children[0]
	checkTextBlock(t, liText, "a")

	liBlock = ulBlock.Children[1]
	checkBlock(t, liBlock, TypeLI, 1)
	liText = liBlock.Children[0]
	checkTextBlock(t, liText, "b")

	liBlock = ulBlock.Children[2]
	checkBlock(t, liBlock, TypeLI, 3)
	liText = liBlock.Children[0]
	checkTextBlock(t, liText, "")
	aBlock := liBlock.Children[1]
	checkAnchorBlock(t, aBlock, "c", "https://mokelab.com")
	liText = liBlock.Children[2]
	checkTextBlock(t, liText, "")
}

func Test_4(t *testing.T) {
	out, err := Parse(src4)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- p
	//    |   |- text
	//    |- pre
	//    |- p
	checkBlock(t, out, TypeRoot, 3)

	pBlock := out.Children[0]
	checkBlock(t, pBlock, TypeP, 1)
	pText := pBlock.Children[0]
	checkTextBlock(t, pText, "Hey")

	preCodeBlock := out.Children[1]
	checkBlock(t, preCodeBlock, TypePreCode, 1)
	codeText := preCodeBlock.Children[0]
	checkTextBlock(t, codeText, "package main\n\nfunc main() {\n}\n")

	pBlock = out.Children[2]
	checkBlock(t, pBlock, TypeP, 1)
	pText = pBlock.Children[0]
	checkTextBlock(t, pText, "text node")

}

func Test_5(t *testing.T) {
	// language is not supported..
	out, err := Parse(src_5)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- p
	//    |  |- text
	//    |- pre
	//       |- text
	checkBlock(t, out, TypeRoot, 2)

	pBlock := out.Children[0]
	checkBlock(t, pBlock, TypeP, 1)

	pText := pBlock.Children[0]
	checkTextBlock(t, pText, "Hey")

	pBlock = out.Children[1]
	checkBlock(t, pBlock, TypePreCode, 1)
	pText = pBlock.Children[0]
	checkTextBlock(t, pText, "package main\n\nfunc main() {\n}\n")
}

func Test_6(t *testing.T) {
	out, err := Parse(src_6)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	checkBlock(t, out, TypeRoot, 2)

	child := out.Children[0]
	checkBlock(t, child, TypeP, 1)

	child = child.Children[0]
	checkTextBlock(t, child, "Hey")

	child = out.Children[1]
	checkBlock(t, child, TypePreCode, 1)
	child = child.Children[0]
	checkTextBlock(t, child, "package main\n\nfunc main() {\n   s := `\n`\n}\n")
}

func Test_7(t *testing.T) {
	out, err := Parse(src7)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	// root
	//    |- h1
	//    |   |- text
	//    |- p
	//       |- text
	//       |- code
	//       |- text
	checkBlock(t, out, TypeRoot, 2)

	h1Block := out.Children[0]
	checkBlock(t, h1Block, TypeH1, 1)

	h1Text := h1Block.Children[0]
	checkTextBlock(t, h1Text, "Program")

	pBlock := out.Children[1]
	checkBlock(t, pBlock, TypeP, 3)

	pText := pBlock.Children[0]
	checkTextBlock(t, pText, "This is ")

	codeBlock := pBlock.Children[1]
	checkInlineCodeBlock(t, codeBlock, "code")

	pText = pBlock.Children[2]
	checkTextBlock(t, pText, " block")
}
