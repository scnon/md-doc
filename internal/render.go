package internal

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func Render2Html(content []byte) string {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(ext)
	doc := parser.Parse(content)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return string(markdown.Render(doc, renderer))
}
