package notifier

import (
	"testing"

	"github.com/ViBiOh/funds/model"
)

func TestInit(t *testing.T) {
	var tests = []struct {
		wantErr error
	}{
		{
			nil,
		},
	}

	for _, test := range tests {
		result := InitEmail()

		if result != test.wantErr {
			t.Errorf("InitEmail() = %v, want %v", result, test.wantErr)
		}
	}
}

func TestGetHTMLContent(t *testing.T) {
	var tests = []struct {
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
			[]byte(`<body style="box-sizing: border-box; width: 100%; border: 0; margin: 0; padding: 5px;"><h1 style="box-sizing: border-box; width: 100%; text-align: center; background-color: #3a3a3a; color: #f8f8f8; border: 0; padding: 5px; margin: 0 5px 0 0;">Funds</h1><p>Bonjour,<p style="color: #4cae4c;">Les fonds suivants viennent de dépasser le score de <strong>0</strong>.<p><table style="box-sizing: border-box; border: 0; margin: 0; padding: 0; width: 100%;"><thead style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="border: 0; margin: 0; padding: 0; width: 100%;"><td style="padding: 5px; width: 140px;">ISIN<td style="padding: 5px;">Libellé<td style="padding: 5px; width: 80px;">Score<tbody style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="background-color: #e1e1e8;"><td style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_1" rel="noopener noreferrer" target=_blank>test_1</a><td style="padding: 5px;">Above Fund<td style="padding: 5px; width: 80px;">80<tr><td style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_2" rel="noopener noreferrer" target=_blank>test_2</a><td style="padding: 5px;">Above Fund Second<td style="padding: 5px; width: 80px;">10</table><p style="color: #d43f3a;">Les fonds suivants viennent de repasser sous leur seuil initial d'alerte.<p><table style="box-sizing: border-box; border: 0; margin: 0; padding: 0; width: 100%;"><thead style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="border: 0; margin: 0; padding: 0; width: 100%;"><td style="padding: 5px; width: 140px;">ISIN<td style="padding: 5px;">Libellé<td style="padding: 5px; width: 80px;">Score<tbody style="border: 0; margin: 0; padding: 0; width: 100%;"><tr style="background-color: #e1e1e8;"><td style="padding: 5px; width: 140px;"><a href="https://funds.vibioh.fr/?isin=test_2" rel="noopener noreferrer" target=_blank>test_2</a><td style="padding: 5px;">Below Fund<td style="padding: 5px; width: 80px;">-10</table><p>Pour plus d'informations, n'hésitez pas à consulter <a href="https://funds.vibioh.fr/?o=score" rel="noopener noreferrer" target=_blank>notre site</a>.<p>Bonne journée,<br>A bientôt,<p>--<br>Funds App - powered by <a href=https://vibioh.fr rel="noopener noreferrer" target=_blank>ViBiOh</a>`),
			nil,
		},
	}

	var failed bool

	for _, test := range tests {
		result, err := getHTMLContent(test.score, test.above, test.below)

		failed = false

		if err == nil && test.wantErr != nil {
			failed = true
		} else if err != nil && test.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != test.wantErr.Error() {
			failed = true
		} else if string(result) != string(test.want) {
			failed = true
		}

		if failed {
			t.Errorf("getHTMLContent(%.2f, %v, %v) = (%s, %v), want (%s, %v)", test.score, test.above, test.below, result, err, test.want, test.wantErr)
		}
	}
}
