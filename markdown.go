package markdown

// Markdown provides API to convert markdown to other language
type Markdown interface {
	// Compile markdown to other language
	Compile(src string) (string, error)
}
