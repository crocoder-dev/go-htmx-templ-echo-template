package main

import (
	"context"
	"go-htmx-templ-echo-template/internals/handlers"
	"log"
	"os"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	e.POST("/create_table_data", app.CreateTableData)
	e.GET("/read_table_data", app.ReadTableData)
	e.PUT("/update_table_data", app.UpdateTableData)
	e.DELETE("/delete_table_data", app.DeleteTableData)

	e.GET("/show_modal", app.ShowModal)
	e.POST("/close_modal", app.CloseModal)

	e.GET("/open_update_row", app.OpenUpdateRow)
	e.GET("/cancel_update", app.CancelUpdate)

	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
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
