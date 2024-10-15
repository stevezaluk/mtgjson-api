package models

type CardLegalities struct {
	Alchemy         string `json:"alchemy"`
	Brawl           string `json:"brawl"`
	Commander       string `json:"commander"`
	Duel            string `json:"duel"`
	Explorer        string `json:"explorer"`
	Future          string `json:"future"`
	Gladiator       string `json:"gladiator"`
	Historic        string `json:"historic"`
	HistoricBrawl   string `json:"historicBrawl"`
	Legacy          string `json:"legacy"`
	Modern          string `json:"modern"`
	OathBreaker     string `json:"oathbreaker"`
	OldSchool       string `json:"oldschool"`
	Pauper          string `json:"pauper"`
	PauperCommander string `json:"paupercommander"`
	Penny           string `json:"penny"`
	Pioneer         string `json:"pioneer"`
	Predh           string `json:"predh"`
	PreModern       string `json:"premodern"`
	Standard        string `json:"standard"`
	StandardBrawl   string `json:"standardbrawl"`
	Timeless        string `json:"timeless"`
	Vintage         string `json:"vintage"`
}
