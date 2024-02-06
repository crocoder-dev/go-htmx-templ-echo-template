package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"go-htmx-templ-echo-template/internals/templates"
	"os"
	"strconv"

	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

var usersData map[int]templates.User

// var id int

var db *sql.DB

func addData() {
	os.Remove("./db/sqlite-database.db")

	file, err := os.Create("./db/sqlite-database.db") // Create SQLite file
	if err != nil {
		fmt.Println(err.Error())
	}
	file.Close()
	db, _ = sql.Open("sqlite3", "./db/sqlite-database.db")
	// defer db.Close()
	createTable(db)

	insertUser(db, "Dean", 28, "New York", "NY")
	insertUser(db, "Sam", 26, "New York", "NY")

	row, err := db.Query("SELECT * FROM user ORDER BY ID")
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()

	usersData = make(map[int]templates.User)

	for row.Next() {
		var ID int
		var Name string
		var Age int
		var City string
		var State string
		row.Scan(&ID, &Name, &Age, &City, &State)

		newUser := templates.User{
			ID:    ID,
			Name:  Name,
			Age:   Age,
			City:  City,
			State: State,
		}

		usersData[ID] = newUser
	}
	// usersDataDummy := make(map[int]templates.User)
	// id = 0
	// usersData = make(map[int]templates.User)

	// usersDataDummy[id] = templates.User{
	// 	ID:    id,
	// 	Name:  "Dean",
	// 	Age:   28,
	// 	City:  "New York",
	// 	State: "NY",
	// }
	// id++
	// usersDataDummy[id] = templates.User{
	// 	ID:    id,
	// 	Name:  "Sam",
	// 	Age:   26,
	// 	City:  "New York",
	// 	State: "NY",
	// }
	// id++

	// keys := make([]int, 0, len(usersDataDummy))

	// for k := range usersDataDummy {
	// 	keys = append(keys, k)
	// }
	// sort.Ints(keys)
	// usersData = usersDataDummy
	// for _, k := range keys {
	// 	usersData[k] = usersDataDummy[k]
	// }
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

	err = insertUser(db, name, ageInt, city, state)
	if err != nil {
		return generateMessage(c, "Error creating user", "error")
	}

	newUser, err := getLastUser(db)
	if err != nil {
		return generateMessage(c, "Error fetching user", "error")
	}

	usersData[newUser.ID] = newUser

	components := templates.UserRow(newUser, true)
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

		err := updateUser(db, user)
		if err != nil {
			return generateMessage(c, "User update error", "error")
		}

		usersData[idInt] = user
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

	if _, ok := usersData[idInt]; ok {
		err := deleteUser(db, idInt)
		if err != nil {
			return generateMessage(c, "User delete error", "error")
		}
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
	if user, ok := usersData[idInt]; ok {
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

func createTable(db *sql.DB) {
	createUserTableSQL := `CREATE TABLE user (
		"ID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Name" TEXT,
		"Age" integer,
		"City" TEXT,	
		"State" TEXT
	  );`

	statement, err := db.Prepare(createUserTableSQL)
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()
}

func getLastUser(db *sql.DB) (templates.User, error) {
	var user templates.User

	query := "SELECT ID, Name, Age, City, State FROM user ORDER BY ID DESC LIMIT 1"

	row := db.QueryRow(query)

	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.State)
	if err != nil {
		return user, err
	}

	return user, nil
}

func insertUser(db *sql.DB, name string, age int, city string, state string) error {
	insertUserSQL := `INSERT INTO user(name, age, city, state) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertUserSQL, name, age, city, state)
	return err
}

func updateUser(db *sql.DB, user templates.User) error {
	updateUserSQL := `UPDATE user SET Name=?, Age=?, City=?, State=? WHERE ID=?`
	_, err := db.Exec(updateUserSQL, user.Name, user.Age, user.City, user.State, user.ID)
	return err
}

func deleteUser(db *sql.DB, id int) error {
	deleteUserSQL := `DELETE FROM user WHERE ID=?`
	_, err := db.Exec(deleteUserSQL, id)
	return err
}
