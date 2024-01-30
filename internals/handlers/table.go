package handlers

import (
	"context"
	"go-htmx-templ-echo-template/internals/templates"
	"sort"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var usersData map[int]templates.User
var id int

func addData() {
	usersDataDummy := make(map[int]templates.User)
	id = 0
	usersData = make(map[int]templates.User)

	usersDataDummy[id] = templates.User{
		ID:    id,
		Name:  "Dean",
		Age:   28,
		City:  "New York",
		State: "NY",
	}
	id++
	usersDataDummy[id] = templates.User{
		ID:    id,
		Name:  "Sam",
		Age:   26,
		City:  "New York",
		State: "NY",
	}
	id++

	keys := make([]int, 0, len(usersDataDummy))

	for k := range usersDataDummy {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	usersData = usersDataDummy
	for _, k := range keys {
		usersData[k] = usersDataDummy[k]
	}
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
		return generateMessage(c, "Invalid Age", "error")
	}

	newUser := templates.User{
		ID:    id,
		Name:  name,
		Age:   ageInt,
		City:  city,
		State: state,
	}
	id++

	usersData[newUser.ID] = newUser

	components := templates.UserRow(newUser, true)
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
		return generateMessage(c, "Invalid Age", "error")
	}
	// Convert ID to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return generateMessage(c, "Invalid ID", "error")
	}

	if user, ok := usersData[idInt]; ok {
		if formName != "" && formName != user.Name {
			user.Name = formName
		}

		if formAgeStr != "" && ageInt != user.Age {
			user.Age = ageInt
		}

		if formCity != "" && formCity != user.City {
			user.City = formCity
		}

		if formState != "" && formState != user.State {
			user.State = formState
		}

		usersData[idInt] = user
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(user, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	// If the user is not found, return 404 Not Found
	return generateMessage(c, "User not found", "error")
}

func (a *App) DeleteRow(c echo.Context) error {
	id := c.QueryParam("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	if _, ok := usersData[idInt]; ok {
		delete(usersData, idInt)
		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
	}
	return generateMessage(c, "User not found", "error")
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
		return generateMessage(c, "Invalid ID", "error")
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

	return generateMessage(c, "User not found", "error")
}

func (a *App) CancelUpdate(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if user, ok := usersData[idInt]; !ok {
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(user, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	return generateMessage(c, "User update cancel failed", "error")
}

func generateMessage(c echo.Context, message string, state string) error {
	// c.Response().WriteHeader(http.StatusNotFound)
	c.Response().Header().Set("HX-Reswap", "beforeend")
	c.Response().Header().Set("HX-Retarget", "#messages")
	components := templates.MessageItem(message, state)
	return components.Render(context.Background(), c.Response().Writer)
}
