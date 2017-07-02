package ast

import "testing"

const src_1 = `
# Hello

World![image](./a.webp)
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

const src_4 = "Hey\n\n" +
	"```\n" +
	"package main\n\n" +
	"func main() {\n" +
	"}\n" +
	"```\n"

func checkTextBlock(t *testing.T, b *Block, value string) {
	if b.Type != TypeText {
		t.Errorf("Type must be text but %d", b.Type)
		return
	}
	if b.Value != value {
		t.Errorf("Value must be %s but %s", value, b.Value)
		return
	}
}

func checkAnchorBlock(t *testing.T, b *Block, text, url string) {
	if b.Type != TypeAnchor {
		t.Errorf("Type must be text but %d", b.Type)
		return
	}
	if b.Value != text {
		t.Errorf("Anchor text must be %s but %s", text, b.Value)
		return
	}
	if b.URL != url {
		t.Errorf("Anchor URL must be %s but %s", url, b.URL)
		return
	}
}

func checkImageBlock(t *testing.T, b *Block, text, url string) {
	if b.Type != TypeImage {
		t.Errorf("Type must be text but %d", b.Type)
		return
	}
	if b.Value != text {
		t.Errorf("Image title must be %s but %s", text, b.Value)
		return
	}
	if b.URL != url {
		t.Errorf("Image URL must be %s but %s", url, b.URL)
		return
	}
}

func checkChildrenCount(t *testing.T, b *Block, value int) {
	if len(b.Children) != value {
		t.Errorf("Children must have %d but %d", value, len(b.Children))
		return
	}
}

func checkBlock(t *testing.T, b *Block, tp BlockType, count int) {
	if b.Type != tp {
		t.Errorf("type must be %d but %d", tp, b.Type)
		return
	}
	if len(b.Children) != count {
		t.Errorf("children count must be %d but %d", count, len(b.Children))
		return
	}
}

func Test_1(t *testing.T) {
	out, err := Parse(src_1)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	checkBlock(t, out, TypeRoot, 2)

	child := out.Children[0]
	checkBlock(t, child, TypeH1, 1)

	child = child.Children[0]
	checkTextBlock(t, child, "Hello")

	child = out.Children[1]
	checkBlock(t, child, TypeP, 1)

	text := child.Children[0]
	checkBlock(t, text, TypeText, 3)
	textChild := text.Children[0]
	checkTextBlock(t, textChild, "World")
	textChild = text.Children[1]
	checkImageBlock(t, textChild, "image", "./a.webp")
	textChild = text.Children[2]
	checkTextBlock(t, textChild, "")
}

func Test_2(t *testing.T) {
	out, err := Parse(src_2)
	if err != nil {
		t.Errorf("Parse error : %s", err)
		return
	}
	checkBlock(t, out, TypeRoot, 2)

	child := out.Children[0]
	checkBlock(t, child, TypeH1, 1)

	child = child.Children[0]
	checkTextBlock(t, child, "Hello")

	child = out.Children[1]
	checkBlock(t, child, TypeH2, 1)

	child = child.Children[0]
	checkTextBlock(t, child, "World")
}

func Test_3(t *testing.T) {
	out, err := Parse(src_3)
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
	checkBlock(t, child, TypeUL, 3)

	li := child.Children[0]
	checkBlock(t, li, TypeLI, 1)
	liText := li.Children[0]
	checkTextBlock(t, liText, "a")

	li = child.Children[1]
	checkBlock(t, li, TypeLI, 1)
	liText = li.Children[0]
	checkTextBlock(t, liText, "b")

	li = child.Children[2]
	checkBlock(t, li, TypeLI, 1)
	liText = li.Children[0]
	checkBlock(t, liText, TypeText, 3)
	text1 := liText.Children[0]
	checkTextBlock(t, text1, "")
	text2 := liText.Children[1]
	checkAnchorBlock(t, text2, "c", "https://mokelab.com")
	text3 := liText.Children[2]
	checkTextBlock(t, text3, "")
}

func Test_4(t *testing.T) {
	out, err := Parse(src_4)
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
	checkTextBlock(t, child, "package main\n\nfunc main() {\n}\n")

}
