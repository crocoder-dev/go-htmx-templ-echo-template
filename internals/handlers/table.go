package handlers

import (
	"context"
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

	ageInt := 0
	fmt.Sscanf(ageStr, "%d", &ageInt)

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
	ageInt := 0
	fmt.Sscanf(formAgeStr, "%d", &ageInt)
	// Convert ID to int
	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

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

func (a *App) DeleteTableData(c echo.Context) error {
	id := c.QueryParam("id")

	fmt.Println("id: ", id)

	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

	for i, item := range TableList {
		if item.ID == idInt {
			TableList = append(TableList[:i], TableList[i+1:]...)
			return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted"})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) ShowModal(c echo.Context) error {
	id := c.QueryParam("id")
	modal_type := c.QueryParam("modal_type")

	action := "/create_table_data"
	method := "POST"
	buttonText := "Create"

	if modal_type == "update" {
		action = "/update_table_data"
		buttonText = "Update"
		method = "PUT"
	}

	item := types.TableItem{
		ID:    0,
		Name:  "",
		Age:   0,
		City:  "",
		State: "",
	}
	idInt := 0
	if id != "" {
		fmt.Sscanf(id, "%d", &idInt)

		for _, row := range TableList {
			if row.ID == idInt {
				item = row
				break
			}
		}
	}

	modalData := types.ModalData{
		Action:     action,
		Data:       item,
		ButtonText: buttonText,
		Method:     method,
	}

	components := templates.Modal(modalData)

	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CloseModal(c echo.Context) error {
	return nil
}
