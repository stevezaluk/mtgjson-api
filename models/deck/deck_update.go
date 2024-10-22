package deck

type DeckUpdate struct {
	MainBoard []string `json:"mainBoard"`
	SideBoard []string `json:"sideBoard"`
	Commander []string `json:"commander"`
}

func (d DeckUpdate) AllCards() []string {
	var ret []string

	ret = append(d.MainBoard, d.SideBoard...)
	ret = append(ret, d.Commander...)

	return ret
}

func (d DeckUpdate) Count() int {
	ret := len(d.MainBoard) + len(d.SideBoard) + len(d.Commander)

	return ret
}
