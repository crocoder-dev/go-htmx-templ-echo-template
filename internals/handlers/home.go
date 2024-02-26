package handlers

import (
	"context"
	"encoding/json"
	"go-htmx-templ-echo-template/internals/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) Home(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	b, _ := json.MarshalIndent(h, "", "\t")

	page := &templates.Page{
		Title:   "HOME",
		Boosted: h.HxBoosted,
	}

	components := templates.Home(page, "Dude!", string(b))
	return components.Render(context.Background(), c.Response().Writer)
}
