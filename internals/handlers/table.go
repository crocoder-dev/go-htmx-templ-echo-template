package handlers

import (
	"context"
	"go-htmx-templ-echo-template/internals/templates"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var tableData map[int]templates.Item
var id int

func addData() {
	id = 0
	tableData = make(map[int]templates.Item)
	tableData[id] = templates.Item{
		ID:    id,
		Name:  "Dean",
		Age:   28,
		City:  "New York",
		State: "NY",
	}
	id++
	tableData[id] = templates.Item{
		ID:    id,
		Name:  "Sam",
		Age:   26,
		City:  "New York",
		State: "NY",
	}
	id++
}

func (a *App) TablePage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}
	addData()

	components := templates.Table(page, tableData, false, nil)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CreateRow(c echo.Context) error {
	name := c.FormValue("name")
	ageStr := c.FormValue("age")
	city := c.FormValue("city")
	state := c.FormValue("state")

	ageInt, err := strconv.Atoi(ageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Age"})
	}

	newItem := templates.Item{
		ID:    id,
		Name:  name,
		Age:   ageInt,
		City:  city,
		State: state,
	}
	id++

	tableData[newItem.ID] = newItem

	components := templates.TableRow(newItem, true)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) UpdateRow(c echo.Context) error {
	// Retrieve form values
	id := c.FormValue("id")
	formName := c.FormValue("name")
	formAgeStr := c.FormValue("age")
	formCity := c.FormValue("city")
	formState := c.FormValue("state")
	// Convert age to int
	ageInt, err := strconv.Atoi(formAgeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Age"})
	}
	// Convert ID to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if item, ok := tableData[idInt]; ok {
		if formName != "" && formName != item.Name {
			item.Name = formName
		}

		if formAgeStr != "" && ageInt != item.Age {
			item.Age = ageInt
		}

		if formCity != "" && formCity != item.City {
			item.City = formCity
		}

		if formState != "" && formState != item.State {
			item.State = formState
		}

		tableData[idInt] = item
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.TableRow(item, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	// If the item is not found, return 404 Not Found
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) DeleteRow(c echo.Context) error {
	id := c.QueryParam("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	if _, ok := tableData[idInt]; ok {
		delete(tableData, idInt)
		return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted"})
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) ShowModal(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}

	var req = c.Request().Header.Get("HX-Request")
	if req != "true" {
		if len(tableData) == 0 {
			addData()
		}

		components := templates.Table(page, tableData, true, nil)
		return components.Render(context.Background(), c.Response().Writer)
	}

	components := templates.Modal()
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CloseModal(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/users")
	return c.NoContent(http.StatusOK)
}

func (a *App) OpenUpdateRow(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	var req = c.Request().Header.Get("HX-Request")
	if req != "true" {
		if len(tableData) == 0 {
			addData()
		}

		page := &templates.Page{
			Title:   "Users demo",
			Boosted: h.HxBoosted,
		}

		components := templates.Table(page, tableData, false, &idInt)
		return components.Render(context.Background(), c.Response().Writer)
	}
	if _, ok := tableData[idInt]; ok {
		components := templates.UsersList(tableData, &idInt)
		return components.Render(context.Background(), c.Response().Writer)
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) CancelUpdate(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if item, ok := tableData[idInt]; ok {
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.TableRow(item, false)
		return components.Render(context.Background(), c.Response().Writer)
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item update cancel failed"})

}
