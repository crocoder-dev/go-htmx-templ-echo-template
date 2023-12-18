package types

import (
	"github.com/donseba/go-htmx"
)

type App struct {
	HTMX *htmx.HTMX
}

type TableItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	City  string `json:"city"`
	State string `json:"state"`
}
