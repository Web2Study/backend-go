package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (
	User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
)

var (
	users = map[int]*User{}
	seq   = 1
)

func createUser(c echo.Context) error {
	u := &User{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	log.Println(u.ID, "--", u.Name)
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}
