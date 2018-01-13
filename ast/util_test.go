package ast

import "testing"

func checkTextBlock(t *testing.T, b *Block, value string) {
	if b.Type != TypeText {
		t.Errorf("Type must be text but %d", b.Type)
		return
	}
	if b.Value != value {
		t.Errorf("Value must be %s but %s len=%d but %d", value, b.Value, len(value), len(b.Value))
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
