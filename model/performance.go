package model

// Performance of a funds
type Performance struct {
	ID            string    `json:"id"`
	Isin          string    `json:"isin"`
	Label         string    `json:"label"`
	Category      string    `json:"category"`
	Rating        string    `json:"rating"`
	OneMonth      float64   `json:"1m"`
	ThreeMonths   float64   `json:"3m"`
	SixMonths     float64   `json:"6m"`
	OneYear       float64   `json:"1y"`
	VolThreeYears float64   `json:"v3y"`
	Score         float64   `json:"score"`
}

// GetID returns Performance's ID
func (p Performance) GetID() string {
	return p.ID
}

// ComputeScore calculate score from Performance's values
func (p *Performance) ComputeScore() {
	score := (0.25 * p.OneMonth) + (0.3 * p.ThreeMonths) + (0.25 * p.SixMonths) + (0.2 * p.OneYear) - (0.1 * p.VolThreeYears)
	p.Score = float64(int(score*100)) / 100
}
