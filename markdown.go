package markdown

// Markdown provides markdown to other API
type Markdown interface {
	// ToHTML compiles markdown to html
	ToHTML(src string) (string, error)
}
