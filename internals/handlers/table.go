package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", "./db/sqlite-database.db")

func (a *App) UsersPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Users demo",
		Boosted: h.HxBoosted,
	}

	if users, err := getAllUsers(); err == nil {
		components := templates.Users(page, users, false, nil)
		return components.Render(context.Background(), c.Response().Writer)
	}
	return c.JSON(http.StatusInternalServerError, "Server Error")
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

	user, err := insertUser(name, ageInt, city, state)
	if err != nil {
		return generateMessage(c, "Error creating user", "error")
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
		return generateMessage(c, "Invalid Age", "error")
	}
	// Convert ID to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return generateMessage(c, "Invalid ID", "error")
	}

	if user, err := getUserByID(idInt); err == nil {
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

		err := updateUser(user)
		if err != nil {
			return generateMessage(c, "User update error", "error")
		}

		c.Response().Header().Set("HX-Push-Url", "/users")
		components := templates.UserRow(user, false)
		return components.Render(context.Background(), c.Response().Writer)
	}

	return generateMessage(c, "User not found", "error")
}

func (a *App) DeleteRow(c echo.Context) error {
	id := c.QueryParam("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	if _, err := getUserByID(idInt); err == nil {
		err := deleteUser(idInt)
		if err != nil {
			return generateMessage(c, "User delete error", "error")
		}
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

	users, err := getAllUsers()
	if err != nil {
		return generateMessage(c, "SQL error", "error")
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
		return generateMessage(c, "Invalid ID", "error")
	}
	users, err := getAllUsers()
	if err != nil {
		return generateMessage(c, "SQL error", "error")
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
	if _, err := getUserByID(idInt); err == nil {
		fmt.Println(err)
		components := templates.UsersList(users, &idInt)
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
	if user, err := getUserByID(idInt); err == nil {
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

func getAllUsers() ([]templates.User, error) {
	rows, err := db.Query("SELECT ID, Name, Age, City, State FROM user ORDER BY ID")
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

func getUserByID(id int) (templates.User, error) {
	getUserSQL := `SELECT ID, Name, Age, City, State FROM user WHERE ID = ?`
	var user templates.User
	err := db.QueryRow(getUserSQL, strconv.Itoa(id)).Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
	if err != nil {
		return user, err
	}
	return user, nil
}

func insertUser(name string, age int, city string, state string) (templates.User, error) {
	insertUserSQL := `INSERT INTO user(name, age, city, state) VALUES (?, ?, ?, ?) RETURNING ID, Name, Age, City, State`
	var user templates.User
	err := db.QueryRow(insertUserSQL, name, age, city, state).Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
	if err != nil {
		return templates.User{}, err
	}

	return user, nil
}

func updateUser(user templates.User) error {
	updateUserSQL := `UPDATE user SET Name=?, Age=?, City=?, State=? WHERE ID=?`
	_, err := db.Exec(updateUserSQL, user.Name, user.Age, user.City, user.State, user.ID)
	return err
}

func deleteUser(id int) error {
	deleteUserSQL := `DELETE FROM user WHERE ID=?`
	_, err := db.Exec(deleteUserSQL, id)
	return err
}
