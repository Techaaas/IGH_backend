package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func main() {
	var db = database{nil}
	db.connector()
	db.dropTables()
	main2()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	//DefaultCORSConfig = CORSConfig{
	//	Skipper:      DefaultSkipper,
	//	AllowOrigins: []string{"*"},
	//	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	//}

	// Routes
	helloHandler := func(c echo.Context) error {
		return hello(c, db)
	}
	getDiffHandler := func(c echo.Context) error {
		return getDiff(c, db)
	}
	getBranchesHandler := func(c echo.Context) error {
		return getBranches(c, db)
	}
	getCommitsHandler := func(c echo.Context) error {
		return getCommits(c, db)
	}
	e.GET("/", helloHandler)
	e.GET("/:branch/commits", getCommitsHandler)
	e.GET("branches", getBranchesHandler)
	e.GET("/diff/:hashFirst/:hashSecond", getDiffHandler)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}

// Handler
func hello(c echo.Context, db database) error {
	return c.JSON(http.StatusOK, "Techaas. All rights reserved. 2023")
}

func getCommits(c echo.Context, db database) error {
	branch := c.Param("branch")
	return c.String(http.StatusOK, strings.Join(db.getAllCommits(branch), "%"))
}

func getBranches(c echo.Context, db database) error {
	return c.String(http.StatusOK, strings.Join(db.getAllBranches(), "%"))
}

func getDiff(c echo.Context, db database) error {
	var commit1 string
	var commit2 string
	commit1 = c.Param("hashFirst")
	commit2 = c.Param("hashSecond")
	var b = db.getDiff(commit1, commit2)
	return c.JSON(http.StatusOK, b)
}
