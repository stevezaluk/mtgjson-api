package deck

type Deck struct {
	Code        string   `json:"code"`
	Commander   []string `json:"commander"`
	MainBoard   []string `json:"mainBoard"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"releaseDate"`
	SideBoard   []string `json:"sideBoard"`
	Type        string   `json:"type"`
}
