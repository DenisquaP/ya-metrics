package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}
