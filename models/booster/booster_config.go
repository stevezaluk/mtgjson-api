package booster

type BoosterConfig struct {
	Boosters    []BoosterPack           `json:"boosters"`
	TotalWeight int64                   `json:"boostersTotalWeight"`
	Name        string                  `json:"name"`
	Sheets      map[string]BoosterSheet `json:"sheets"`
}
