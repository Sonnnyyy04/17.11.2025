package models

type Batch struct {
	ID       int               `json:"id"`
	Links    []string          `json:"links"`
	Statuses map[string]string `json:"statuses"`
}

