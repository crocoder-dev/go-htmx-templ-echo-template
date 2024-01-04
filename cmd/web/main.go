package main

import (
	"context"
	"go-htmx-templ-echo-template/internals/handlers"

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

	e.GET("/table", app.Table)
	e.POST("/table", app.CreateTableData)
	e.PUT("/table", app.UpdateTableData)
	e.DELETE("/table", app.DeleteTableData)
	e.PATCH("/table", app.Table)

	e.GET("/table/new", app.ShowModal)

	e.GET("/table/update/:slug", app.OpenUpdateRow)
	e.POST("/table/update/:slug", app.CancelUpdate)

	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":3000"))
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
