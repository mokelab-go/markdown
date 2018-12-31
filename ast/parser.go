package ast

import (
	"errors"
	"fmt"
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

	stateReadBlock = iota
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
	textLinkStart
	textLinkEnd
	textLinkUrlStart
	text_link_url_end
	textImage1
	textImageStart
	textImageEnd
	textImageURLStart
	textImageURLAttrKey
	textImageURLAttrValue
	text_image_url_end
	text_br
	textCodeStart
	textCodeEnd
	textInlineCode
)

func Parse(src string) (*Block, error) {
	var textValue []byte
	var attrKey []byte
	var attrValue []byte
	state := stateReadBlock
	index := 0
	srcLen := len(src)
	root := newBlock(TypeRoot)
	currentBlock := root
	blockStack := &blockStack{values: make([]*Block, 0)}

	for index < srcLen {
		char := src[index]
		switch state {
		case stateReadBlock:
			if char == '#' {
				state = read_h1
			} else if char == '`' {
				state = read_bq_1
			} else if char == '\n' {
				if currentBlock.Type == TypeLI {
					// pop x2
					blockStack.Pop()
					blockStack.Pop()
				}
				// ignore
			} else if char == ' ' {
				state = read_head_space
			} else {
				// Add p-text block
				pBlock := newBlock(TypeP)
				textBlock := newBlock(TypeText)
				appendChild(currentBlock, pBlock)
				appendChild(pBlock, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(pBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
				index--
			}
		case read_bq_1:
			if char == '`' {
				state = read_bq_2
			} else {
				// textnode + start code block
				pBlock := newBlock(TypeP)
				codeBlock := newBlock(TypeCode)
				appendChild(currentBlock, pBlock)
				appendChild(pBlock, codeBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(pBlock)

				state = textInlineCode
				currentBlock = codeBlock
				textValue = make([]byte, 0)
				index--
			}
		case read_bq_2:
			if char == '`' {
				state = read_bq_3
			} else {
				pBlock := newBlock(TypeP)
				textBlock := newBlock(TypeText)
				appendChild(currentBlock, pBlock)
				appendChild(pBlock, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(pBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
				textValue = appendStr(textValue, "``")
				index--
			}
		case read_bq_3:
			if char == '\n' {

				codeBlock := newBlock(TypePreCode)
				textBlock := newBlock(TypeText)
				appendChild(currentBlock, codeBlock)
				appendChild(codeBlock, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(codeBlock)

				state = code
				currentBlock = textBlock
				textValue = make([]byte, 0)
			} else {
				// TODO : check lang
				pBlock := newBlock(TypeP)
				textBlock := newBlock(TypeText)
				appendChild(currentBlock, pBlock)
				appendChild(pBlock, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(pBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
				textValue = appendStr(textValue, "```")

				index--
			}
		case code:
			if char == '`' {
				state = read_bq_end_1
			} else if char == '<' {
				textValue = appendStr(textValue, "&lt;")
			} else if char == '>' {
				textValue = appendStr(textValue, "&gt;")
			} else {
				textValue = append(textValue, char)
			}
		case read_bq_end_1:
			if char == '`' {
				state = read_bq_end_2
			} else {
				textValue = append(textValue, '`')
				textValue = append(textValue, char)
				state = code
			}
		case read_bq_end_2:
			if char == '`' {
				state = read_bq_end_3
			} else {
				textValue = appendStr(textValue, "``")
				textValue = append(textValue, char)
				state = code
			}
		case read_bq_end_3:
			if char == '\n' {
				currentBlock.Value = string(textValue)

				state = stateReadBlock
				currentBlock = blockStack.Pop()
				currentBlock = blockStack.Pop()
			} else if char == ' ' {
				// ignore
			} else {
				return nil, errors.New(fmt.Sprintf("\\n expected byte %d at %d", char, index))
			}
		case read_head_space:
			if char == '*' || char == '-' {
				state = read_ul
			} else {
				state = text
			}
		case read_ul:
			if char == ' ' {
				if currentBlock.Type == TypeUL {
					liBlock := newBlock(TypeLI)
					textBlock := newBlock(TypeText)
					appendChild(currentBlock, liBlock)
					appendChild(liBlock, textBlock)
					blockStack.Push(currentBlock)
					blockStack.Push(liBlock)

					state = text
					currentBlock = textBlock
					textValue = make([]byte, 0)
				} else {
					ulBlock := newBlock(TypeUL)
					liBlock := newBlock(TypeLI)
					textBlock := newBlock(TypeText)
					appendChild(currentBlock, ulBlock)
					appendChild(ulBlock, liBlock)
					appendChild(liBlock, textBlock)
					blockStack.Push(currentBlock)
					blockStack.Push(ulBlock)
					blockStack.Push(liBlock)

					state = text
					currentBlock = textBlock
					textValue = make([]byte, 0)
				}
			} else {
				return nil, errors.New(fmt.Sprintf("unexpected token at %d", index))
			}
		case text:
			if char == '\n' {
				parentBlock := blockStack.Top()
				if parentBlock.Type == TypeLI {
					currentBlock.Value = string(textValue)

					state = stateReadBlock
					// pop x 2
					currentBlock = blockStack.Pop()
					currentBlock = blockStack.Pop()
				} else {
					state = text_br
				}
			} else if char == '[' {
				currentBlock.Value = string(textValue)

				// add link node
				aBlock := newBlock(TypeAnchor)
				appendChild(blockStack.Top(), aBlock)

				state = textLinkStart
				currentBlock = aBlock
				textValue = make([]byte, 0)
			} else if char == '!' {
				state = textImage1
			} else if char == '`' {
				currentBlock.Value = string(textValue)

				// add inline code node
				inlineCodeBlock := newBlock(TypeCode)
				appendChild(blockStack.Top(), inlineCodeBlock)

				state = textInlineCode
				currentBlock = inlineCodeBlock
				textValue = make([]byte, 0)
			} else if char == '<' {
				// nop
			} else if char == '>' {
				// nop
			} else {
				textValue = append(textValue, char)
			}
		case textLinkStart:
			if char == ']' {
				currentBlock.Value = string(textValue)

				state = textLinkEnd
				textValue = make([]byte, 0)
			} else {
				textValue = append(textValue, char)
			}
		case textLinkEnd:
			if char == '(' {
				state = textLinkUrlStart
			} else {
				return nil, fmt.Errorf("expected is ( but %d at %d", char, index)
			}
		case textLinkUrlStart:
			if char == ')' {
				currentBlock.URL = string(textValue)

				// add next text block
				textBlock := newBlock(TypeText)
				appendChild(blockStack.Top(), textBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
			} else {
				textValue = append(textValue, char)
			}
		case textImage1:
			if char == '[' {
				currentBlock.Value = string(textValue)

				imageBlock := newBlock(TypeImage)
				appendChild(blockStack.Top(), imageBlock)

				state = textImageStart
				currentBlock = imageBlock
				textValue = make([]byte, 0)
			} else {
				textValue = append(textValue, '!')
				index--
				state = text
			}
		case textImageStart:
			if char == ']' {
				currentBlock.Value = string(textValue)

				state = textImageEnd
				textValue = make([]byte, 0)
			} else {
				textValue = append(textValue, char)
			}
		case textImageEnd:
			if char == '(' {
				state = textImageURLStart
			} else {
				return nil, fmt.Errorf("expected is ( but %d at %d", char, index)
			}
		case textImageURLStart:
			if char == ')' {
				currentBlock.URL = string(textValue)

				// set next text block
				textBlock := newBlock(TypeText)
				appendChild(blockStack.Top(), textBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
			} else if char == ' ' {
				currentBlock.URL = string(textValue)
				currentBlock.Attributes = make(map[string]string)

				state = textImageURLAttrKey
				attrKey = make([]byte, 0)
			} else {
				textValue = append(textValue, char)
			}
		case textImageURLAttrKey:
			if char == '=' {
				state = textImageURLAttrValue
				attrValue = make([]byte, 0)
			} else {
				attrKey = append(attrKey, char)
			}
		case textImageURLAttrValue:
			if char == ' ' {
				currentBlock.Attributes[string(attrKey)] = string(attrValue)

				state = textImageURLAttrKey
				attrKey = make([]byte, 0)
			} else if char == ')' {
				currentBlock.Attributes[string(attrKey)] = string(attrValue)

				// set next text block
				textBlock := newBlock(TypeText)
				appendChild(blockStack.Top(), textBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
			} else {
				attrValue = append(attrValue, char)
			}
		case text_br:
			if char == '\n' {
				// end of text and parent block
				currentBlock.Value = string(textValue)
				// pop x 2
				currentBlock = blockStack.Pop()
				currentBlock = blockStack.Pop()

				state = stateReadBlock
			} else {
				state = text
				textValue = append(textValue, ' ')
				index--
			}
		case read_h1: // #
			if char == ' ' {
				h1Block := newBlock(TypeH1)
				textBlock := newBlock(TypeText)

				appendChild(currentBlock, h1Block)
				appendChild(h1Block, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(h1Block)

				currentBlock = textBlock
				textValue = make([]byte, 0)

				state = text
			} else if char == '#' {
				state = read_h2
			}
		case read_h2: // ##
			if char == ' ' {
				h2Block := newBlock(TypeH2)
				textBlock := newBlock(TypeText)
				appendChild(currentBlock, h2Block)
				appendChild(h2Block, textBlock)
				blockStack.Push(currentBlock)
				blockStack.Push(h2Block)

				currentBlock = textBlock
				textValue = make([]byte, 0)

				state = text
			}
		case textInlineCode:
			if char == '`' {
				currentBlock.Value = string(textValue)

				// set next text block
				textBlock := newBlock(TypeText)
				appendChild(blockStack.Top(), textBlock)

				state = text
				currentBlock = textBlock
				textValue = make([]byte, 0)
			} else if char == '<' {
				textValue = appendStr(textValue, "&lt;")
			} else if char == '>' {
				textValue = appendStr(textValue, "&gt;")
			} else {
				textValue = append(textValue, char)
			}
		}
		index++
	}
	if currentBlock.Type == TypeText {
		currentBlock.Value = string(textValue)
	}

	return root, nil
}

func appendChild(b *Block, c *Block) {
	b.Children = append(b.Children, c)
}

func appendStr(out []byte, text string) []byte {
	return append(out, text...)
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
