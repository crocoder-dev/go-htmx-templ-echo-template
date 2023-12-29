package handlers

import (
	"context"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"go-htmx-templ-echo-template/internals/types"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var TableList []types.TableItem

func (a *App) Table(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	TableList = TableList[:0]
	TableList = append(TableList, types.TableItem{
		ID:    1,
		Name:  "Dude",
		Age:   25,
		City:  "New York",
		State: "NY",
	})

	page := &templates.Page{
		Title:   "Home",
		Boosted: h.HxBoosted,
	}

	components := templates.Table(page, TableList)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CreateTableData(c echo.Context) error {
	name := c.FormValue("name")
	ageStr := c.FormValue("age")
	city := c.FormValue("city")
	state := c.FormValue("state")

	ageInt, err := strconv.Atoi(ageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid age")
	}

	newItem := types.TableItem{
		ID:    0,
		Name:  name,
		Age:   ageInt,
		City:  city,
		State: state,
	}

	if len(TableList) == 0 {
		newItem.ID = 0
		TableList = []types.TableItem{newItem}
	} else {
		newItem.ID = TableList[len(TableList)-1].ID + 1
		TableList = append(TableList, newItem)
	}

	components := templates.TableRow(newItem)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) ReadTableData(c echo.Context) error {
	return c.JSON(http.StatusOK, TableList)
}

func (a *App) UpdateTableData(c echo.Context) error {
	// Retrieve form values
	id := c.QueryParam("id")
	formName := c.FormValue("name")
	formAgeStr := c.FormValue("age")
	formCity := c.FormValue("city")
	formState := c.FormValue("state")
	// Convert age to int
	ageInt, err := strconv.Atoi(formAgeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid age")
	}
	// Convert ID to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	for i, item := range TableList {
		if item.ID == idInt {
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

			TableList[i] = item
			components := templates.TableRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}

	// If the item is not found, return 404 Not Found
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) OpenUpdateRow(c echo.Context) error {
	id := c.QueryParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	for _, item := range TableList {
		if item.ID == idInt {
			components := templates.TableInputRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) CancelUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	for _, item := range TableList {
		if item.ID == idInt {
			components := templates.TableRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item update cancel failed"})
}

func (a *App) DeleteTableData(c echo.Context) error {
	id := c.QueryParam("id")

	fmt.Println("id: ", id)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	for i, item := range TableList {
		if item.ID == idInt {
			TableList = append(TableList[:i], TableList[i+1:]...)
			return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted"})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) ShowModal(c echo.Context) error {
	components := templates.Modal()
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CloseModal(c echo.Context) error {
	return nil
}
