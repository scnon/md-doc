package utils

import (
	"bytes"
	"text/template"

	"github.com/scnon/md-doc/model"
)

func RenderSearchItem(data []model.SearchItem) string {
	tmpl, err := template.ParseFiles("./static/search_item.html")
	if err != nil {
		return ""
	}

	var reader bytes.Buffer
	tmpl.Execute(&reader, map[string]interface{}{
		"Data": data,
		"Last": len(data) - 1,
	})

	return reader.String()
}
