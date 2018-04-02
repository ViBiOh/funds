package notifier

import (
	"testing"

	"github.com/ViBiOh/funds/pkg/model"
)

func TestInit(t *testing.T) {
	var cases = []struct {
		wantErr error
	}{
		{
			nil,
		},
	}

	for _, testCase := range cases {
		result := InitEmail()

		if result != testCase.wantErr {
			t.Errorf(`InitEmail() = %v, want %v`, result, testCase.wantErr)
		}
	}
}

func TestGetHTMLContent(t *testing.T) {
	var cases = []struct {
		score   float64
		above   []*model.Fund
		below   []*model.Fund
		want    []byte
		wantErr error
	}{
		{
			0.0,
			nil,
			nil,
			nil,
			nil,
		},
		{
			0.0,
			[]*model.Fund{{Isin: `test_1`, Score: 80.00, Label: `Above Fund`}, {Isin: `test_2`, Score: 10.00, Label: `Above Fund Second`}},
			[]*model.Fund{{Isin: `test_2`, Score: -10.00, Label: `Below Fund`}},
			[]byte(`<body class="no-space padding full-width" style="border: 0; margin: 0; padding: 5px; width: 100%;"><h1 class="no-space padding full-width center" style="background-color: #3a3a3a; border: 0; color: #f8f8f8; margin: 0 5px 0 0; padding: 5px; text-align: center; width: 100%;" align=center>Funds</h1><p>Bonjour,<p style="color: #4cae4c;">Les fonds suivants viennent de dépasser le score de <strong>0</strong>.<p><table class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><thead class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><tr class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><td class="padding isin" style="padding: 5px; width: 140px;">ISIN<td class=padding style="padding: 5px;">Libellé<td class="padding score" style="padding: 5px; width: 80px;">Score<tbody class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="background-color: #e1e1e8;"><td class="padding isin" style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_1" rel="noopener noreferrer" target=_blank>test_1</a><td class=padding style="padding: 5px;">Above Fund<td class="padding score" style="padding: 5px; width: 80px;">80<tr><td class="padding isin" style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_2" rel="noopener noreferrer" target=_blank>test_2</a><td class=padding style="padding: 5px;">Above Fund Second<td class="padding score" style="padding: 5px; width: 80px;">10</table><p style="color: #d43f3a;">Les fonds suivants viennent de repasser sous leur seuil initial d&#39;alerte.<p><table class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><thead class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><tr class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><td class="padding isin" style="padding: 5px; width: 140px;">ISIN<td class=padding style="padding: 5px;">Libellé<td class="padding score" style="padding: 5px; width: 80px;">Score<tbody class="no-space full-width" style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="background-color: #e1e1e8;"><td class="padding isin" style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_2" rel="noopener noreferrer" target=_blank>test_2</a><td class=padding style="padding: 5px;">Below Fund<td class="padding score" style="padding: 5px; width: 80px;">-10</table><p>Pour plus d&#39;informations, n&#39;hésitez pas à consulter <a href="https://funds.vibioh.fr/?o=score" rel="noopener noreferrer" target=_blank>notre site</a>.<p>Bonne journée,<br>A bientôt,<p>--<br>Funds App - powered by <a href=https://vibioh.fr rel="noopener noreferrer" target=_blank>ViBiOh</a>`),
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := getHTMLContent(testCase.score, testCase.above, testCase.below)

		failed = false

		if err == nil && testCase.wantErr != nil {
			failed = true
		} else if err != nil && testCase.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != testCase.wantErr.Error() {
			failed = true
		} else if string(result) != string(testCase.want) {
			failed = true
		}

		if failed {
			t.Errorf(`getHTMLContent(%.2f, %v, %v) = (%s, %v), want (%s, %v)`, testCase.score, testCase.above, testCase.below, result, err, testCase.want, testCase.wantErr)
		}
	}
}
