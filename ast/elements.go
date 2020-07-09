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
	TypeCode
	TypeText
	TypeAnchor
	TypeImage
)

type Block struct {
	Type       BlockType
	URL        string
	Value      string
	Children   []*Block
	Attributes map[string]string
}

func newBlock(t BlockType) *Block {
	return &Block{
		Type:       t,
		Value:      "",
		Children:   make([]*Block, 0),
		Attributes: make(map[string]string),
	}
}

type blockStack struct {
	values []*Block
}

func (s *blockStack) Push(v *Block) {
	s.values = append(s.values, v)
}

func (s *blockStack) Pop() *Block {
	top := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return top
}

func (s *blockStack) Top() *Block {
	return s.values[len(s.values)-1]
}

func (s *blockStack) Clear() {
	s.values = make([]*Block, 0)
}

func appendChild(b *Block, c *Block) {
	b.Children = append(b.Children, c)
}

func appendStr(out []byte, text string) []byte {
	return append(out, text...)
}
