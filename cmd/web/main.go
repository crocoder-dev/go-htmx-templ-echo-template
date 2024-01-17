package main

import (
	"context"
	"go-htmx-templ-echo-template/internals/handlers"
	"go-htmx-templ-echo-template/internals/templates"
	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {

	app := &handlers.App{
		HTMX: htmx.New(),
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(HtmxMiddleware)

	e.GET("/", app.Home)
	e.GET("/about", app.About)

	e.GET("/users", app.TablePage)

	e.POST("/users", app.CreateRow)
	e.PUT("/users", app.UpdateRow)
	e.DELETE("/users", app.DeleteRow)

	e.DELETE("/users/new", app.CloseModal)
	e.GET("/users/new", app.ShowModal)

	e.GET("/users/update/:id", app.OpenUpdateRow)
	e.POST("/users/update/:id", app.CancelUpdate)

	e.HTTPErrorHandler = HTTPErrorHandler

	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":3000"))
}

func HTTPErrorHandler(err error, c echo.Context) {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	page := &templates.Page{
		Title:   "404 Not Found",
		Boosted: h.HxBoosted,
	}

	if code == http.StatusNotFound {
		components := templates.NotFound(page)
		if err := components.Render(context.Background(), c.Response().Writer); err != nil {
			c.Logger().Error(err)
		}
	}
}

func HtmxMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxRequestHeader{
			HxBoosted:               htmx.HxStrToBool(c.Request().Header.Get("HX-Boosted")),
			HxCurrentURL:            c.Request().Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(c.Request().Header.Get("HX-History-Restore-Request")),
			HxPrompt:                c.Request().Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(c.Request().Header.Get("HX-Request")),
			HxTarget:                c.Request().Header.Get("HX-Target"),
			HxTriggerName:           c.Request().Header.Get("HX-Trigger-Name"),
			HxTrigger:               c.Request().Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
