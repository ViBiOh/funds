package notifier

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ViBiOh/funds/model"
	"github.com/aymerick/douceur/inliner"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type scoreTemplateContent struct {
	Score      float64
	AboveFunds []*model.Fund
	BelowFunds []*model.Fund
}

var mailTmpl *template.Template
var minifier *minify.M

// InitEmail initialize template and minifier
func InitEmail() error {
	mailTmpl = template.Must(template.New(`email.html`).Funcs(template.FuncMap{
		`odd`: func(i int) bool {
			return i%2 == 0
		},
	}).ParseFiles(`email.html`))

	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	minifier = m

	return nil
}

func getHTMLContent(scoreLevel float64, above []*model.Fund, below []*model.Fund) ([]byte, error) {
	if len(above) == 0 && len(below) == 0 {
		return nil, nil
	}

	templateBuffer := &bytes.Buffer{}
	if err := mailTmpl.ExecuteTemplate(templateBuffer, `main`, scoreTemplateContent{Score: scoreLevel, AboveFunds: above, BelowFunds: below}); err != nil {
		return nil, fmt.Errorf(`Error while executing template: %v`, err)
	}

	html, err := inliner.Inline(templateBuffer.String())
	if err != nil {
		return nil, fmt.Errorf(`Error while inlining style: %v`, err)
	}

	templateBuffer.Reset()
	if _, err = templateBuffer.WriteString(html); err != nil {
		return nil, fmt.Errorf(`Error while loading buffer: %v`, err)
	}

	minifyBuffer := &bytes.Buffer{}
	if err := minifier.Minify(`text/html`, minifyBuffer, templateBuffer); err != nil {
		return nil, fmt.Errorf(`Error while minifying template: %v`, err)
	}

	return minifyBuffer.Bytes(), nil
}
