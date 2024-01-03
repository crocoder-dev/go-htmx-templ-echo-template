package handlers

import (
	"context"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var tableData map[int]templates.Item
var id int

func (a *App) Table(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	fmt.Println("hx-boosted:", h.HxBoosted)
	fmt.Println("hx-current-url:", h.HxCurrentURL)
	fmt.Println("hx-history-restore-request:", h.HxHistoryRestoreRequest)
	fmt.Println("hx-prompt:", h.HxPrompt)
	fmt.Println("hx-target:", h.HxTarget)
	fmt.Println("hx-trigger:", h.HxTrigger)
	fmt.Println("hx-trigger-name:", h.HxTriggerName)

	tableData = make(map[int]templates.Item)
	id = 0
	tableData[id] = templates.Item{
		ID:    id,
		Name:  "Dean",
		Age:   28,
		City:  "New York",
		State: "NY",
	}
	id++

	page := &templates.Page{
		Title:   "Home",
		Boosted: h.HxBoosted,
	}

	components := templates.Table(page, tableData)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CreateTableData(c echo.Context) error {
	name := c.FormValue("name")
	ageStr := c.FormValue("age")
	city := c.FormValue("city")
	state := c.FormValue("state")

	ageInt, err := strconv.Atoi(ageStr)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "<p>Invalid age</p>")
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

	components := templates.TableRow(newItem)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) UpdateTableData(c echo.Context) error {
	// Retrieve form values
	id := c.FormValue("id")
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

	for i, item := range tableData {
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

			tableData[i] = item
			components := templates.TableRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}

	// If the item is not found, return 404 Not Found
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) OpenUpdateRow(c echo.Context) error {
	id := c.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	for _, item := range tableData {
		if item.ID == idInt {
			components := templates.TableInputRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) CancelUpdate(c echo.Context) error {
	id := c.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	for _, item := range tableData {
		if item.ID == idInt {
			components := templates.TableRow(item)
			return components.Render(context.Background(), c.Response().Writer)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item update cancel failed"})
}

func (a *App) DeleteTableData(c echo.Context) error {
	id := c.FormValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	for _, item := range tableData {
		if item.ID == idInt {
			delete(tableData, idInt)
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
