package notifier

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ViBiOh/funds/model"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

const scoreNotificationTemplate = `
{{ define "main" }}
<body style="box-sizing: border-box; width: 100%; border: 0; margin: 0; padding: 5px;">
	<h1 style="box-sizing: border-box; width: 100%; text-align: center; background-color: #3a3a3a; color: #f8f8f8; border: 0; padding: 5px; margin: 0 5px 0 0;">Funds</h1>
	<p>Bonjour,</p>
	{{ if len .AboveFunds }}
		<p style="color: #4cae4c;">Les fonds suivants viennent de dépasser le score de <strong>{{ .Score }}</strong>.</p>
		<p>
			{{ template "funds" .AboveFunds }}
		</p>
	{{ end }}
	{{ if len .BelowFunds }}
		<p style="color: #d43f3a;">Les fonds suivants viennent de repasser sous leur seuil initial d'alerte.</p>
		<p>
			{{ template "funds" .BelowFunds }}
		</p>
	{{ end }}
	<p>
		Pour plus d'informations, n'hésitez pas à consulter <a href="https://funds.vibioh.fr/?o=score" rel="noopener noreferrer" target="_blank">notre site</a>.
	</p>
	<p>
		Bonne journée,
		<br />
		A bientôt,
	</p>
	<p>
		--
		<br />
		Funds App - powered by <a href="https://vibioh.fr" rel="noopener noreferrer" target="_blank">ViBiOh</a>
	</p>
</body>
{{ end }}

{{ define "funds" }}
<table style="box-sizing: border-box; border: 0; margin: 0; padding: 0; width: 100%;">
	<thead style="border: 0; margin: 0; padding: 0; width: 100%;">
		<tr style="border: 0; margin: 0; padding: 0; width: 100%;">
			<td style="padding: 5px; width: 140px;">ISIN</td>
			<td style="padding: 5px;">Libellé</td>
			<td style="padding: 5px; width: 80px;">Score</td>
		</tr>
	</thead>
	<tbody style="border: 0; margin: 0; padding: 0; width: 100%;">
		{{ range $index, $fund := . }}
			<tr style="{{ if odd $index }}background-color: #e1e1e8;{{ end }}">
				<td style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin={{ $fund.Isin }}" rel="noopener noreferrer" target="_blank">{{ $fund.Isin }}</a></td>
				<td style="padding: 5px;">{{ $fund.Label }}</td>
				<td style="padding: 5px; width: 80px;">{{ $fund.Score }}</td>
			</tr>
		{{end}}
	</tbody>
</table>
{{ end }}
`

type scoreTemplateContent struct {
	Score      float64
	AboveFunds []*model.Fund
	BelowFunds []*model.Fund
}

var mailTmpl *template.Template
var minifier *minify.M

// InitEmail initialize template and minifier
func InitEmail() error {
	tmpl := template.New(`ScoreNotification`)

	tmpl.Funcs(template.FuncMap{`odd`: func(i int) bool {
		return i%2 == 0
	}})

	tmpl, err := tmpl.Parse(scoreNotificationTemplate)
	if err != nil {
		return fmt.Errorf(`Error while parsing template: %v`, err)
	}

	mailTmpl = tmpl

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

	minifyBuffer := &bytes.Buffer{}
	if err := minifier.Minify("text/html", minifyBuffer, templateBuffer); err != nil {
		return nil, fmt.Errorf(`Error while minifying template: %v`, err)
	}

	return minifyBuffer.Bytes(), nil
}
