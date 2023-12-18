package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"go-htmx-templ-echo-template/internals/types"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var TableList []types.TableItem

func (a *App) Table(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

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

	components := templates.Table(page, "Dude!", TableList)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CreateTableData(c echo.Context) error {
	var newItem types.TableItem
	err := json.NewDecoder(c.Request().Body).Decode(&newItem)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	newItem.ID = len(TableList) + 1

	TableList = append(TableList, newItem)
	return c.JSON(http.StatusCreated, newItem)
}

func (a *App) ReadTableData(c echo.Context) error {
	return c.JSON(http.StatusOK, TableList)
}

func (a *App) UpdateTableData(c echo.Context) error {
	var partialUpdate map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&partialUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Get the ID from the request parameters
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID parameter is missing"})
	}

	// Convert the ID to an integer
	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

	// Find the item by ID
	for i, item := range TableList {
		if item.ID == idInt {
			// Update only non-empty fields provided in the request
			if name, ok := partialUpdate["Name"].(string); ok && name != "" {
				item.Name = name
			}

			if age, ok := partialUpdate["Age"].(float64); ok {
				item.Age = int(age)
			}

			if city, ok := partialUpdate["City"].(string); ok && city != "" {
				item.City = city
			}

			if state, ok := partialUpdate["State"].(string); ok && state != "" {
				item.State = state
			}

			// Add more cases for other fields as needed

			// Update the item in the list
			TableList[i] = item
			return c.JSON(http.StatusOK, item)
		}
	}

	// If the item is not found, return 404 Not Found
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) DeleteTableData(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID parameter is missing"})
	}

	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

	for i, item := range TableList {
		if item.ID == idInt {
			TableList = append(TableList[:i], TableList[i+1:]...)
			return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted successfully"})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}
