package notifier

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/ViBiOh/funds/model"
)

const scoreNotificationTemplate = `
{{ define "main" }}
<body style="border: 0; margin: 0; padding: 5px;">
	<h1 style="width: 100%; text-align: center; background-color: #3a3a3a; color: #f8f8f8; border: 0; padding: 5px; margin: 0;">Funds</h1>
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
		Pour plus d'information, n'hésitez pas à consulter <a href="https://funds.vibioh.fr/?o=score" rel="noopener noreferrer" target="_blank">notre site</a>.
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
<table style="border: 0; margin: 0; padding: 0; width: 100%;">
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
	AboveFunds []model.Performance
	BelowFunds []model.Performance
}

var mailTmpl *template.Template

func init() {
	tmpl := template.New(`ScoreNotification`)

	tmpl.Funcs(template.FuncMap{`odd`: func(i int) bool {
		return i%2 == 0
	}})

	tmpl, err := tmpl.Parse(scoreNotificationTemplate)
	if err != nil {
		log.Fatal(err)
	}

	mailTmpl = tmpl
}

func getHTMLContent(scoreLevel float64, above []model.Performance, below []model.Performance) ([]byte, error) {
	buffer := &bytes.Buffer{}

	if err := mailTmpl.ExecuteTemplate(buffer, `main`, scoreTemplateContent{Score: scoreLevel, AboveFunds: above, BelowFunds: below}); err != nil {
		return nil, fmt.Errorf(`Error while creating HTML content: %v`, err)
	}

	return buffer.Bytes(), nil
}
