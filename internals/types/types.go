package types

type TableItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	City  string `json:"city"`
	State string `json:"state"`
}
