package model

// Fund informations
type Fund struct {
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
func (f Fund) GetID() string {
	return f.ID
}

// ComputeScore calculate score of Fund
func (f *Fund) ComputeScore() {
	score := (0.25 * f.OneMonth) + (0.3 * f.ThreeMonths) + (0.25 * f.SixMonths) + (0.2 * f.OneYear) - (0.1 * f.VolThreeYears)
	f.Score = float64(int(score*100)) / 100
}
