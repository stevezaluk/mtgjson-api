package deck

type DeckUpdate struct {
	MainBoard []string `json:"mainBoard"`
	SideBoard []string `json:"sideBoard"`
	Commander []string `json:"commander"`
}
