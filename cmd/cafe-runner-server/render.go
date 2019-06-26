package main

import (
	"html/template"
	"io"

	"strings"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
)

func parseTemplates(box *packr.Box) (*template.Template, error) {
	root := template.New("")

	err := box.Walk(func(path string, file packr.File) error {
		if strings.HasSuffix(path, ".html") {
			path = strings.Replace(path, `\`, `/`, -1)
			b := file.String()

			t := root.New(path)
			t, e2 := t.Parse(b)
			if e2 != nil {
				return e2
			}

		}
		return nil
	})
	return root, err
}

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
