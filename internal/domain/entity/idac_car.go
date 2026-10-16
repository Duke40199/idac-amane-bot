package entity

// CarSpecInfo holds the canonical model code and base spec for lookups
type CarSpecInfo struct {
	Maker         string
	CarName       string
	ModelCode     string // Canonical model code (e.g. "CZ4A")
	BaseSpec      string // "speed" or "tech"
	SpecStyleName string
	SegaSpecID    string
	Aliases       []string
}

type CarMetadata struct {
	ID                  string  `json:"id"` // pk
	SegaCarID           int64   `json:"car_id"`
	Name                string  `json:"car_name"`
	ModelCode           string  `json:"model_code"`
	Maker               string  `json:"maker_name"`
	BaseStyleName       string  `json:"base_style_name"`
	NormalizedBaseStyle string  `json:"normalized_base_style"`
	CarStyleIDs         []int64 `json:"style_car_id"`
	// Aggregated fields
	SpecIDs   []string `json:"spec_ids"`
	SpecNames []string `json:"spec_names"`
	// for searching
	SearchBlob string `json:"-"`
}

type CarStyleMetadata struct {
	ID             string `json:"id"` // pk
	StyleCarID     int64  `json:"style_car_id"`
	CarID          string `json:"car_id"` // fk of cars
	Name           string `json:"style_name"`
	RouteStyleName string `json:"route_style_name"`
}
