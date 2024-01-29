package handlers

import (
	"context"
	"go-htmx-templ-echo-template/internals/templates"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var usersData map[int]templates.User
var id int

func addData() {
	id = 0
	usersData = make(map[int]templates.User)
	usersData[id] = templates.User{
		ID:    id,
		Name:  "Dean",
		Age:   28,
		City:  "New York",
		State: "NY",
	}
	id++
	usersData[id] = templates.User{
		ID:    id,
		Name:  "Sam",
		Age:   26,
		City:  "New York",
		State: "NY",
	}
	id++
}

func (a *App) UsersPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}
	addData()

	components := templates.Users(page, usersData, false, nil)
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

	newItem := templates.User{
		ID:    id,
		Name:  name,
		Age:   ageInt,
		City:  city,
		State: state,
	}
	id++

	usersData[newItem.ID] = newItem

	components := templates.UserRow(newItem, true)
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

	if item, ok := usersData[idInt]; ok {
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

		usersData[idInt] = item
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(item, false)
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

	if _, ok := usersData[idInt]; ok {
		delete(usersData, idInt)
		return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted"})
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func (a *App) ShowAddUserModal(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}

	var req = c.Request().Header.Get("HX-Request")
	if req != "true" {
		if len(usersData) == 0 {
			addData()
		}

		components := templates.Users(page, usersData, true, nil)
		return components.Render(context.Background(), c.Response().Writer)
	}

	components := templates.Modal()
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CloseAddUserModal(c echo.Context) error {
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
		if len(usersData) == 0 {
			addData()
		}

		page := &templates.Page{
			Title:   "Users demo",
			Boosted: h.HxBoosted,
		}

		components := templates.Users(page, usersData, false, &idInt)
		return components.Render(context.Background(), c.Response().Writer)
	}
	if _, ok := usersData[idInt]; ok {
		components := templates.UsersList(usersData, &idInt)
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
	if item, ok := usersData[idInt]; ok {
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(item, false)
		return components.Render(context.Background(), c.Response().Writer)
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item update cancel failed"})
}
