package ast

import (
	"errors"
)

type stateFunc func(s *parseState, char byte) (stateFunc, error)

type parseState struct {
	src          string
	index        int
	srcLen       int
	root         *Block
	currentBlock *Block
	blockStack   *blockStack

	textValue      []byte
	linkTitleValue []byte
	linkURLValue   []byte
	attrName       []byte
	attrValue      []byte

	hCount int
}

// Parse src markdown to block
func Parse(src string) (*Block, error) {
	root := newBlock(TypeRoot)
	s := &parseState{
		src:          src,
		index:        0,
		srcLen:       len(src),
		root:         root,
		currentBlock: root,
		blockStack:   &blockStack{values: make([]*Block, 0)},
		hCount:       0,
	}
	f := stateReadRootBlock
	panicCounter := 0
	for s.index < s.srcLen {
		panicCounter++
		if panicCounter > s.srcLen*10 {
			return nil, errors.New("parser may be in infinte loop")
		}
		char := src[s.index]

		f2, err := f(s, char)
		if err != nil {
			return nil, err
		}

		f = f2
	}
	if s.currentBlock.Type == TypeText {
		s.currentBlock.Value = string(s.textValue)
	}
	return s.root, nil
}

func stateReadRootBlock(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' || char == '\n' {
		// skip
		s.index++
		return stateReadRootBlock, nil
	}
	if char == '#' {
		s.hCount = 1
		s.index++
		return stateReadHn, nil
	}
	if char == '*' || char == '-' {
		s.index++
		return stateReadUL, nil
	}
	if char == '`' {
		s.hCount = 1
		s.index++
		return stateReadBeginPreCode, nil
	}
	// paragraph block
	pBlock := newBlock(TypeP)
	textBlock := newBlock(TypeText)
	appendChild(s.currentBlock, pBlock)
	appendChild(pBlock, textBlock)
	s.blockStack.Push(s.currentBlock)
	s.blockStack.Push(pBlock)

	s.currentBlock = textBlock
	s.textValue = make([]byte, 0)
	// read this char as a part of text
	return stateReadText, nil
}

func stateReadHn(s *parseState, char byte) (stateFunc, error) {
	if char == '#' {
		if s.hCount > 2 {
			return nil, errors.New("Cannot support level 2 header")
		}
		s.hCount++
		s.index++
		return stateReadHn, nil
	}
	hBlock := newBlock(toHnType(s.hCount))
	textBlock := newBlock(TypeText)

	appendChild(s.currentBlock, hBlock)
	appendChild(hBlock, textBlock)
	s.blockStack.Push(s.currentBlock)
	s.blockStack.Push(hBlock)
	s.currentBlock = textBlock

	s.textValue = make([]byte, 0)

	if char == ' ' {
		return stateFindFirstText, nil
	}

	return stateReadText, nil
}

func stateFindFirstText(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		s.index++
		return stateFindFirstText, nil
	}
	return stateReadText, nil
}

func stateReadText(s *parseState, char byte) (stateFunc, error) {
	if char == '\n' {
		parentBlock := s.blockStack.Top()
		if parentBlock.Type == TypeH1 ||
			parentBlock.Type == TypeH2 {
			s.currentBlock.Value = string(s.textValue)
			parentBlock = s.blockStack.Pop() // h block is ended
			parentBlock = s.blockStack.Pop() // parent of h block
			s.currentBlock = parentBlock

			s.index++
			return stateReadRootBlock, nil
		}
		if parentBlock.Type == TypeLI {
			s.currentBlock.Value = string(s.textValue)
			parentBlock = s.blockStack.Pop() // li block is ended
			parentBlock = s.blockStack.Pop() // parent of li
			s.currentBlock = parentBlock
			s.index++
			return stateReadNextLiToken, nil
		}
		// ignore \n, if \n again, we must close current block
		s.index++
		return stateReadTextNewLine, nil
	}
	if char == '[' {
		s.linkTitleValue = make([]byte, 0)
		s.index++
		return stateReadLinkTitle, nil
	}
	if char == '!' {
		s.index++
		return stateReadBeginImageToken, nil
	}
	if char == '`' {
		s.currentBlock.Value = string(s.textValue)

		parentBlock := s.blockStack.Top()
		codeBlock := newBlock(TypeCode)
		appendChild(parentBlock, codeBlock)
		s.currentBlock = codeBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadInlineCode, nil
	}
	s.textValue = append(s.textValue, char)
	s.index++
	return stateReadText, nil
}

func stateReadTextNewLine(s *parseState, char byte) (stateFunc, error) {
	if char == '\n' {
		// close all block
		s.currentBlock.Value = string(s.textValue)
		s.blockStack.Clear()
		s.currentBlock = s.root
		s.index++
		return stateReadRootBlock, nil
	}
	s.textValue = append(s.textValue, char)
	s.index++
	return stateReadText, nil
}

// link

func stateReadLinkTitle(s *parseState, char byte) (stateFunc, error) {
	if char == ']' {
		s.index++
		return stateReadLinkURLBeginToken, nil
	}
	s.linkTitleValue = append(s.linkTitleValue, char)
	s.index++
	return stateReadLinkTitle, nil
}

// stateReadLinkURLBeginToken wants '('
func stateReadLinkURLBeginToken(s *parseState, char byte) (stateFunc, error) {
	if char == '(' {
		s.linkURLValue = make([]byte, 0)
		s.index++
		return stateReadLinkURL, nil
	}
	s.textValue = append(s.textValue, '[')
	s.textValue = appendStr(s.textValue, string(s.linkTitleValue))
	s.textValue = append(s.textValue, ']')
	// read current char as a part of text
	return stateReadText, nil
}

func stateReadLinkURL(s *parseState, char byte) (stateFunc, error) {
	if char == ')' {
		s.currentBlock.Value = string(s.textValue)

		parentBlock := s.blockStack.Top()
		linkBlock := newBlock(TypeAnchor)
		linkBlock.Value = string(s.linkTitleValue)
		linkBlock.URL = string(s.linkURLValue)
		appendChild(parentBlock, linkBlock)

		// next block
		textBlock := newBlock(TypeText)
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	s.linkURLValue = append(s.linkURLValue, char)
	s.index++
	return stateReadLinkURL, nil
}

// image

// stateReadBeginImageToken wants '['
func stateReadBeginImageToken(s *parseState, char byte) (stateFunc, error) {
	if char == '[' {
		s.linkTitleValue = make([]byte, 0)
		s.index++
		return stateReadImageTitle, nil
	}
	s.textValue = append(s.textValue, '!')
	// read this character as a part of text
	return stateReadText, nil
}

func stateReadImageTitle(s *parseState, char byte) (stateFunc, error) {
	if char == ']' {
		s.index++
		return stateReadImageURLBeginToken, nil
	}
	s.linkTitleValue = append(s.linkTitleValue, char)
	s.index++
	return stateReadImageTitle, nil
}

// stateReadImageURLBeginToken wants '('
func stateReadImageURLBeginToken(s *parseState, char byte) (stateFunc, error) {
	if char == '(' {
		s.linkURLValue = make([]byte, 0)
		s.index++
		return stateReadImageURL, nil
	}
	s.textValue = append(s.textValue, '!')
	s.textValue = append(s.textValue, '[')
	s.textValue = appendStr(s.textValue, string(s.linkTitleValue))
	s.textValue = append(s.textValue, ']')
	// read current char as a part of text
	return stateReadText, nil
}

func stateReadImageURL(s *parseState, char byte) (stateFunc, error) {
	if char == ')' {
		s.currentBlock.Value = string(s.textValue)

		parentBlock := s.blockStack.Top()
		imageBlock := newBlock(TypeImage)
		imageBlock.Value = string(s.linkTitleValue)
		imageBlock.URL = string(s.linkURLValue)
		appendChild(parentBlock, imageBlock)

		// next block
		textBlock := newBlock(TypeText)
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	if char == ' ' {
		s.currentBlock.Value = string(s.textValue)

		parentBlock := s.blockStack.Top()
		imageBlock := newBlock(TypeImage)
		imageBlock.Value = string(s.linkTitleValue)
		imageBlock.URL = string(s.linkURLValue)
		appendChild(parentBlock, imageBlock)

		s.currentBlock = imageBlock

		s.index++
		return stateReadBeginImageAttr, nil
	}
	s.linkURLValue = append(s.linkURLValue, char)
	s.index++
	return stateReadImageURL, nil
}

func stateReadBeginImageAttr(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		s.index++
		return stateReadBeginImageAttr, nil
	}
	if char == ')' {
		parentBlock := s.blockStack.Top()
		// next block
		textBlock := newBlock(TypeText)
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	s.attrName = make([]byte, 0)
	return stateReadImageAttrName, nil
}

func stateReadImageAttrName(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		// empty value attr
		attrName := string(s.attrName)
		attrValue := ""
		s.currentBlock.Attributes[attrName] = attrValue
		s.index++
		return stateReadBeginImageAttr, nil
	}
	if char == ')' {
		// empty value attr
		attrName := string(s.attrName)
		attrValue := ""
		s.currentBlock.Attributes[attrName] = attrValue

		parentBlock := s.blockStack.Top()
		// next block
		textBlock := newBlock(TypeText)
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	if char == '=' {
		s.attrValue = make([]byte, 0)
		s.index++
		return stateReadImageAttrValue, nil
	}
	s.attrName = append(s.attrName, char)
	s.index++
	return stateReadImageAttrName, nil
}

func stateReadImageAttrValue(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		attrName := string(s.attrName)
		attrValue := string(s.attrValue)
		s.currentBlock.Attributes[attrName] = attrValue
		s.index++
		return stateReadBeginImageAttr, nil
	}
	if char == ')' {
		attrName := string(s.attrName)
		attrValue := string(s.attrValue)
		s.currentBlock.Attributes[attrName] = attrValue

		parentBlock := s.blockStack.Top()
		// next block
		textBlock := newBlock(TypeText)
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	s.attrValue = append(s.attrValue, char)
	s.index++
	return stateReadImageAttrValue, nil
}

// UL

func stateReadUL(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		if s.currentBlock.Type == TypeRoot {
			// put ul
			ulBlock := newBlock(TypeUL)
			appendChild(s.currentBlock, ulBlock)
			s.blockStack.Push(s.currentBlock)
			s.currentBlock = ulBlock
			s.index++
			return stateReadFirstLiToken, nil
		}
		if s.currentBlock.Type == TypeUL {
			s.index++
			return stateReadFirstLiToken, nil
		}
	}
	// just a paragraph
	pBlock := newBlock(TypeP)
	textBlock := newBlock(TypeText)
	appendChild(s.currentBlock, pBlock)
	appendChild(pBlock, textBlock)
	s.blockStack.Push(s.currentBlock)
	s.blockStack.Push(pBlock)

	s.currentBlock = textBlock
	s.textValue = make([]byte, 0)
	s.textValue = append(s.textValue, '*')
	// read this char as a part of text
	return stateReadText, nil
}

func stateReadFirstLiToken(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		// skip
		s.index++
		return stateReadFirstLiToken, nil
	}
	// make li
	liBlock := newBlock(TypeLI)
	textBlock := newBlock(TypeText)
	appendChild(s.currentBlock, liBlock)
	appendChild(liBlock, textBlock)
	s.blockStack.Push(s.currentBlock)
	s.blockStack.Push(liBlock)

	s.currentBlock = textBlock
	s.textValue = make([]byte, 0)
	// read this char as a part of text
	return stateReadText, nil
}

func stateReadNextLiToken(s *parseState, char byte) (stateFunc, error) {
	if char == ' ' {
		// skip
		s.index++
		return stateReadNextLiToken, nil
	}
	if char == '*' || char == '-' {
		s.index++
		return stateReadFirstLiToken, nil
	}
	if char == '\n' {
		// ul block is ended. clear all stack
		s.blockStack.Clear()
		s.currentBlock = s.root
		s.index++
		return stateReadRootBlock, nil
	}
	// new paragraph begins
	s.blockStack.Clear()
	s.currentBlock = s.root

	pBlock := newBlock(TypeP)
	textBlock := newBlock(TypeText)
	appendChild(s.currentBlock, pBlock)
	appendChild(pBlock, textBlock)
	s.blockStack.Push(s.currentBlock)
	s.blockStack.Push(pBlock)

	s.currentBlock = textBlock
	s.textValue = make([]byte, 0)
	// read this char as a part of text
	return stateReadText, nil
}

// pre code
func stateReadBeginPreCode(s *parseState, char byte) (stateFunc, error) {
	if char == '`' {
		s.hCount++
		if s.hCount < 3 {
			s.index++
			return stateReadBeginPreCode, nil
		}
		s.index++
		return stateReadBeginPreCodeNewLine, nil
	}
	if s.hCount == 1 {
		s.index++
		return stateReadText, nil
	}
	s.textValue = make([]byte, 0)
	return stateReadText, nil
}

func stateReadBeginPreCodeNewLine(s *parseState, char byte) (stateFunc, error) {
	if char == '\n' {
		// begin pre code
		preCodeBlock := newBlock(TypePreCode)
		textBlock := newBlock(TypeText)
		appendChild(s.currentBlock, preCodeBlock)
		appendChild(preCodeBlock, textBlock)
		s.blockStack.Push(s.currentBlock)
		s.blockStack.Push(preCodeBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadPreCodeText, nil
	}
	// TODO support language
	s.index++
	return stateReadBeginPreCodeNewLine, nil
}

func stateReadPreCodeText(s *parseState, char byte) (stateFunc, error) {
	if char == '`' {
		s.hCount = 1
		s.index++
		return stateReadEndPreCode, nil
	}
	s.textValue = append(s.textValue, char)
	s.index++
	return stateReadPreCodeText, nil
}

func stateReadEndPreCode(s *parseState, char byte) (stateFunc, error) {
	if char == '`' {
		s.hCount++
		if s.hCount < 3 {
			s.index++
			return stateReadEndPreCode, nil
		}
		// end of pre code
		s.currentBlock.Value = string(s.textValue)

		parentBlock := s.blockStack.Pop() // preCode
		parentBlock = s.blockStack.Pop()  // parent of preCode
		s.currentBlock = parentBlock
		s.index++
		return stateReadRootBlock, nil
	}
	for i := 0; i < s.hCount; i++ {
		s.textValue = append(s.textValue, '`')
	}
	return stateReadPreCodeText, nil
}

// inline code

func stateReadInlineCode(s *parseState, char byte) (stateFunc, error) {
	if char == '`' {
		s.currentBlock.Value = string(s.textValue)

		textBlock := newBlock(TypeText)
		parentBlock := s.blockStack.Top()
		appendChild(parentBlock, textBlock)
		s.currentBlock = textBlock

		s.textValue = make([]byte, 0)
		s.index++
		return stateReadText, nil
	}
	s.textValue = append(s.textValue, char)
	s.index++
	return stateReadInlineCode, nil
}

func toHnType(level int) BlockType {
	switch level {
	case 1:
		return TypeH1
	case 2:
		return TypeH2
	default:
		return TypeH1
	}
}
