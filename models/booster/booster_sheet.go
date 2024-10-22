package booster

type BoosterSheet struct {
	AllowDuplicates bool           `json:"allowDuplicates"`
	BalanceColors   bool           `json:"balanceColors"`
	Cards           map[string]int `json:"cards"`
	Foil            bool           `json:"foil"`
	Fixed           bool           `json:"fixed"`
	TotalWeight     int64          `json:"totalWeight"`
}
