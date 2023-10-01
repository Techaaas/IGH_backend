package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	var db = database{nil}
	db.connector()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	helloHandler := func(c echo.Context) error {
		return hello(c, db)
	}
	e.GET("/", helloHandler)
	e.GET("/bye", goodbye)
	e.GET("/:branch/commits", getCommits)
	e.GET("branches", getBranches)
	e.GET("/diff/:hashFirst/:hashSecond", getDiff)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}

// Handler
func hello(c echo.Context, db database) error {
	var a = db.getDiff("fcb3c93728e025ebcd2e1390c7df716d86cf7644", "dc89ec1b5724cb6fc0d1eb3a9e2910a74a9a43d5")
	fmt.Println(a + " IM HERE")
	return c.JSON(http.StatusOK, a)
}

func goodbye(c echo.Context) error {
	return c.String(http.StatusOK, "Bye bye!!")
}

func getCommits(c echo.Context) error {
	branch := c.Param("branch")
	// DO SOMETHING
	return c.String(http.StatusOK, branch)
}

func getBranches(c echo.Context) error {
	return c.String(http.StatusOK, "List of all branches: [хахахахахаххахахах]")
}

func getDiff(c echo.Context) error {
	commit1 := c.Param("hashFirst")
	commit2 := c.Param("hashSecond")
	return c.String(http.StatusOK, commit1+" diff "+commit2)
}
