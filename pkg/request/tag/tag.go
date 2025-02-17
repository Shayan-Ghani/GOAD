package tagrequest

type Base struct {
	ItemID string `json:"item_id"`
}

type BasePayload struct {
	ItemID string   `json:"item_id"`
	Tags   []string `json:"tags"`
}

type Delete struct {
	Name string `json:"name"`
}

type Tag struct {
	Tags []string `json:"tags"`
}