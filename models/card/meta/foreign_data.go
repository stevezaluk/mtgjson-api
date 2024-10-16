package card

type ForeignData struct {
	FaceName    string          `json:"faceName"`
	FlavorText  string          `json:"flavorText"`
	Identifiers CardIdentifiers `json:"identifiers"`
	Language    string          `json:"language"`
	Name        string          `json:"name"`
	Text        string          `json:"text"`
	Type        string          `json:"type"`
}
