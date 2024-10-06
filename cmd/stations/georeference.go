package stations

type Georeference struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}