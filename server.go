package main

import (
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
	getDiffh := func(c echo.Context) error {
		return getDiff(c, db)
	}
	e.GET("/", helloHandler)
	e.GET("/bye", goodbye)
	e.GET("/:branch/commits", getCommits)
	e.GET("branches", getBranches)
	e.GET("/diff/:hashFirst/:hashSecond", getDiffh)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}

// Handler
func hello(c echo.Context, db database) error {
	//fmt.Println(a + " IM HERE")
	return c.JSON(http.StatusOK, "")
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

func getDiff(c echo.Context, db database) error {
	var commit1 string
	var commit2 string
	commit1 = c.Param("hashFirst")
	commit2 = c.Param("hashSecond")
	var b = db.getDiff(commit1, commit2)
	return c.String(http.StatusOK, b)
}
