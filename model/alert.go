package model

// Alert for a funds
type Alert struct {
	Isin      string
	AlertType string
	Score     float64
}
