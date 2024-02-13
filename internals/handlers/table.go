package handlers

import (
	"context"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func (a *App) UsersPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}

	users, err := a.getAllUsers()

	if err != nil {
		return err
	}

	components := templates.Users(page, users, false, nil)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) CreateRow(c echo.Context) error {
	name := c.FormValue("name")
	ageStr := c.FormValue("age")
	city := c.FormValue("city")
	state := c.FormValue("state")

	ageInt, err := strconv.Atoi(ageStr)
	if err != nil {
		return generateMessage(c, "Invalid Age", "error", http.StatusBadRequest)
	}

	user, err := a.insertUser(name, ageInt, city, state)
	if err != nil {
		return generateMessage(c, "Error creating user", "error", http.StatusInternalServerError)
	}

	components := templates.UserRow(user, true)
	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) UpdateRow(c echo.Context) error {
	id := c.FormValue("id")
	formName := c.FormValue("name")
	formAgeStr := c.FormValue("age")
	formCity := c.FormValue("city")
	formState := c.FormValue("state")
	// Convert age to int
	ageInt, err := strconv.Atoi(formAgeStr)
	if err != nil {
		return generateMessage(c, "Invalid Age", "error", http.StatusBadRequest)
	}
	// Convert ID to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return generateMessage(c, "Invalid ID", "error", http.StatusBadRequest)
	}

	if user, err := a.getUserByID(idInt); err == nil {
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

		err := a.updateUser(user)
		if err != nil {
			return generateMessage(c, "User update error", "error", http.StatusInternalServerError)
		}

		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(user, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	return generateMessage(c, "User not found", "error", http.StatusNotFound)
}

func (a *App) DeleteRow(c echo.Context) error {
	id := c.QueryParam("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return generateMessage(c, "Invalid ID", "error", http.StatusBadRequest)
	}

	if _, err := a.getUserByID(idInt); err == nil {
		err := a.deleteUser(idInt)
		if err != nil {
			return generateMessage(c, "User delete error", "error", http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
	}
	return generateMessage(c, "User not found", "error", http.StatusNotFound)
}

func (a *App) ShowAddUserModal(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}

	users, err := a.getAllUsers()
	if err != nil {
		return generateMessage(c, "SQL error", "error", http.StatusInternalServerError)
	}

	var req = c.Request().Header.Get("HX-Request")
	if req != "true" {
		components := templates.Users(page, users, true, nil)
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
		return generateMessage(c, "Invalid ID", "error", http.StatusBadRequest)
	}
	users, err := a.getAllUsers()
	if err != nil {
		return generateMessage(c, "SQL error", "error", http.StatusBadRequest)
	}

	var req = c.Request().Header.Get("HX-Request")
	if req != "true" {
		page := &templates.Page{
			Title:   "Users demo",
			Boosted: h.HxBoosted,
		}

		components := templates.Users(page, users, false, &idInt)
		return components.Render(context.Background(), c.Response().Writer)

	}
	if _, err := a.getUserByID(idInt); err == nil {
		fmt.Println(err)
		components := templates.UsersList(users, &idInt)
		return components.Render(context.Background(), c.Response().Writer)
	}

	return generateMessage(c, "User not found", "error", http.StatusNotFound)
}

func (a *App) CancelUpdate(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if user, err := a.getUserByID(idInt); err == nil {
		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(user, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	return generateMessage(c, "User update cancel failed", "error", http.StatusInternalServerError)
}

func generateMessage(c echo.Context, message string, state string, status int) error {
	return c.HTML(status, generateMessageHtml(message, state))
}

func (a *App) getAllUsers() ([]templates.User, error) {
	rows, err := a.DB.Query("SELECT id, name, age, city, state FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []templates.User
	for rows.Next() {
		var user templates.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (a *App) getUserByID(id int) (templates.User, error) {
	getUserSQL := `SELECT id, name, age, city, state FROM users WHERE id = ?`
	var user templates.User
	err := a.DB.QueryRow(getUserSQL, strconv.Itoa(id)).Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *App) insertUser(name string, age int, city string, state string) (templates.User, error) {
	insertUserSQL := `INSERT INTO users(name, age, city, state) VALUES (?, ?, ?, ?) RETURNING id, name, age, city, state`
	var user templates.User
	err := a.DB.QueryRow(insertUserSQL, name, age, city, state).Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
	if err != nil {
		return templates.User{}, err
	}

	return user, nil
}

func (a *App) updateUser(user templates.User) error {
	updateUserSQL := `UPDATE users SET name=?, age=?, city=?, state=? WHERE id=?`
	_, err := a.DB.Exec(updateUserSQL, user.Name, user.Age, user.City, user.State, user.ID)
	return err
}

func (a *App) deleteUser(id int) error {
	deleteUserSQL := `DELETE FROM users WHERE id=?`
	_, err := a.DB.Exec(deleteUserSQL, id)
	return err
}

func generateMessageHtml(message string, state string) string {
	if state == "error" {
		return "<div class='mt-2 px-3 py-1 rounded-lg text-white bg-red-500 text-white p-2 rounded-lg'>" + message + "</div>"
	} else {
		return "<div class='mt-2 px-3 py-1 rounded-lg text-white bg-green-500 text-white p-2 rounded-lg'>" + message + "</div>"
	}
}
