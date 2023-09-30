package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/bye", goodbye)
	e.GET("/:branch/commits", getCommits)
	e.GET("branches", getBranches)
	e.GET("/diff/:hashFirst/:hashSecond", getDiff)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}

// Handler
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, "{\n  \"squadName\" : \"Super Hero Squad\",\n  \"homeTown\" : \"Metro City\",\n  \"formed\" : 2016,\n  \"secretBase\" : \"Super tower\",\n  \"active\" : true,\n  \"members\" : [\n    {\n      \"name\" : \"Molecule Man\",\n      \"age\" : 29,\n      \"secretIdentity\" : \"Dan Jukes\",\n      \"powers\" : [\n        \"Radiation resistance\",\n        \"Turning tiny\",\n        \"Radiation blast\"\n      ]\n    },\n    {\n      \"name\" : \"Madame Uppercut\",\n      \"age\" : 39,\n      \"secretIdentity\" : \"Jane Wilson\",\n      \"powers\" : [\n        \"Million tonne punch\",\n        \"Damage resistance\",\n        \"Superhuman reflexes\"\n      ]\n    },\n    {\n      \"name\" : \"Eternal Flame\",\n      \"age\" : 1000000,\n      \"secretIdentity\" : \"Unknown\",\n      \"powers\" : [\n        \"Immortality\",\n        \"Heat Immunity\",\n        \"Inferno\",\n        \"Teleportation\",\n        \"Interdimensional travel\"\n      ]\n    }\n  ]\n}\n")
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
