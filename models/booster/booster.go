package booster

type BoosterPack struct {
	Contents map[string]int `json:"contents"`
	Weight   int64          `json:"weight"`
}
