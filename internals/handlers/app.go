package handlers

import (
	"github.com/donseba/go-htmx"
	"database/sql"
)

type App struct {
	HTMX *htmx.HTMX
	DB   *sql.DB
}
