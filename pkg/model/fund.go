package model

import "context"

// Fund informations
type Fund struct {
	Alert         *Alert  `json:"alert,omitempty"`
	ID            string  `json:"id"`
	Isin          string  `json:"isin"`
	Label         string  `json:"label"`
	Category      string  `json:"category"`
	Rating        string  `json:"rating"`
	OneMonth      float64 `json:"1m"`
	ThreeMonths   float64 `json:"3m"`
	SixMonths     float64 `json:"6m"`
	OneYear       float64 `json:"1y"`
	VolThreeYears float64 `json:"v3y"`
	Score         float64 `json:"score"`
}

// GetID returns Fund's ID
func (f *Fund) GetID() string {
	return f.ID
}

// ComputeScore calculate score of Fund
func (f *Fund) ComputeScore() {
	score := (0.25 * f.OneMonth) + (0.3 * f.ThreeMonths) + (0.25 * f.SixMonths) + (0.2 * f.OneYear) - (0.1 * f.VolThreeYears)
	f.Score = float64(int(score*100)) / 100
}

// GetFundsAbove retrieves funds above score
func (a *App) GetFundsAbove(ctx context.Context, score float64, currentAlerts map[string]Alert) ([]Fund, error) {
	var fundsToAlert []Fund

	funds, err := a.listFundsWithScoreAbove(ctx, score)
	if err != nil {
		return nil, err
	}

	for _, fund := range funds {
		if alert, ok := currentAlerts[fund.Isin]; ok {
			if alert.AlertType != "above" {
				fundsToAlert = append(fundsToAlert, fund)
			}
		} else {
			fundsToAlert = append(fundsToAlert, fund)
		}
	}

	return fundsToAlert, nil
}

// GetFundsBelow retrieves funds below score
func (a *App) GetFundsBelow(ctx context.Context, currentAlerts map[string]Alert) ([]Fund, error) {
	var funds []Fund

	for _, alert := range currentAlerts {
		fund, err := a.readFundByIsin(ctx, alert.Isin)
		if err != nil {
			return nil, err
		}

		if fund.Score < alert.Score {
			funds = append(funds, fund)
		}
	}

	return funds, nil
}
