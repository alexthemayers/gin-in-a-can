package render

import (
	"bytes"
	"fmt"
	"github.com/sasswart/gin-in-a-can/config"
	"github.com/sasswart/gin-in-a-can/sanitizer"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"text/template"
)

// Render contains the main parsing and rendering steps within the render library
// It wraps text/template and creates the output to be written to io.Writer
func Render(config config.Config, data any, templateFile string) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	templater := template.New(templateFile)

	templater.Funcs(templateFuncMap)

	parsedTemplate, err := templater.ParseGlob(fmt.Sprintf("%s/*.tmpl", config.TemplateDirectory))
	if err != nil {
		return nil, err
	}

	err = parsedTemplate.Execute(buff, data)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

var templateFuncMap = template.FuncMap{
	"ToUpper":  strings.ToUpper,
	"ToTitle":  toTitleCase,
	"Sanitize": sanitizer.GoSanitize,
}

func toTitleCase(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}