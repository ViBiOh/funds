package notifier

import (
	"bytes"
	"html/template"

	"github.com/ViBiOh/funds/model"
)

const fundsTemplate = `
<table style="border: 0; margin: 0; padding: 0; width: 100%;">
	<thead style="border: 0; margin: 0; padding: 0; width: 100%;>
		<tr style="border: 0; margin: 0; padding: 0; width: 100%;>
			<td style="padding: 5px;">ISIN</td>
			<td style="padding: 5px;">Libellé</td>
			<td style="padding: 5px;">Score</td>
		</tr>
	</thead>
	<tbody style="border: 0; margin: 0; padding: 0; width: 100%;>
		{{range $index,$fund := .}}
			<tr style="{{if odd $index}}background-color: #e1e1e8;{{end}}">
				<td style="padding: 5px;"><a href="https://funds.vibioh.fr/?isin={{$fund.Isin}}" rel="noopener noreferrer" target="_blank">{{$fund.Isin}}</a></td>
				<td style="padding: 5px;">{{$fund.Label}}</td>
				<td style="padding: 5px;">{{$fund.Score}}</td>
			</tr>
		{{end}}
	</tbody>
</table>
`

const scoreTemplate = `
<body style="border: 0; margin: 0; padding: 5px;">
	<h1 style="width: 100%; text-align: center; background-color: #3a3a3a; color: #f8f8f8; border: 0; padding: 5px; margin: 0;">Funds</h1>
	<p>Bonjour,</p>
	<p>Les fonds suivants ont un score venant de dépasser <strong>{{.Score}}</strong>.</p>
	<p>
		{{template "funds" .Funds}}
	</p>
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
`

type scoreTemplateContent struct {
	Score float64
	Funds []model.Performance
}

func getHTMLContent(scoreLevel float64, funds []model.Performance) ([]byte, error) {
	buffer := &bytes.Buffer{}

	tmpl := template.New(`score`)

	tmpl.Funcs(template.FuncMap{`odd`: func(i int) bool {
		return i%2 == 0
	}})

	tmpl, err := tmpl.Parse(scoreTemplate)
	if err != nil {
		return nil, err
	}

	tmpl, err = tmpl.New(`funds`).Parse(fundsTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buffer, scoreTemplateContent{Score: scoreLevel, Funds: funds}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
