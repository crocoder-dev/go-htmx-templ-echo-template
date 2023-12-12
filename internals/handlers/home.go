package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"go-htmx-templ-echo-template/internals/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) Hello(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	fmt.Println("hello, sign in", c.Get("isSignedIn"))
	fmt.Println("user", c.Get("clerk_user"))

	b, _ := json.MarshalIndent(h, "", "\t")

	page := &templates.Page{
		Title:   "OMG",
		Boosted: h.HxBoosted,
	}

	components := templates.Hello(page, "David!", string(b))
	return components.Render(context.Background(), c.Response().Writer)
}
