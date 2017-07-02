package ast

type BlockType int

const (
	TypeRoot BlockType = iota + 1
	TypeP
	TypeH1
	TypeH2
	TypeUL
	TypeLI
	TypePreCode
	TypeText
	TypeAnchor
	TypeImage
)

type Block struct {
	Type     BlockType
	URL      string
	Value    string
	Children []*Block
}

func newBlock(t BlockType) *Block {
	return &Block{
		Type:     t,
		Value:    "",
		Children: make([]*Block, 0),
	}
}
