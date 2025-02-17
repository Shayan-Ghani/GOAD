package itemrequest

import "github.com/Shayan-Ghani/GOAD/internal/model"

type ID struct {
	ID string
}

type BasePayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DueDate     string   `json:"due_date"`
	Tags        []string `json:"tags"`
}

type Add struct {
	BasePayload
}

type Get struct {
	ID string
}

type GetByTag struct {
	Tags []string
}

type GetDone struct {
	ID  string
	All bool
}

type GetResponse struct {
	Err   string       `json:"error,omitempty"`
	Items []model.Item `json:"items"`
}

type Update struct {
	ID string
	BasePayload
}

type UpdateStatus struct {
	ID string
}

type Delete struct {
	ID string
}
