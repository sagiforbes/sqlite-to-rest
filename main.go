package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/sagiforbes/sqlite-to-rest/utils"
)

var dbFile = flag.String("f", "", "sqlite db file")
var port = flag.String("p", ":4080", "port to listen to. default 4080")

func getTable(c echo.Context) error {
	sql := fmt.Sprintf("SELECT * FROM %s", c.Param("table"))
	res, err := utils.DbQuery(*dbFile, sql)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func main() {
	flag.Parse()
	if _, err := os.Stat(*dbFile); err != nil {
		log.Panic("Invalid sqlite file ", err)
	}

	if err := utils.DbCheckFile(*dbFile); err != nil {
		log.Panic("Invalid sqlite file ", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Pre(middleware.RemoveTrailingSlash())

	//------- register endpoints
	e.GET("/:table", getTable)

	fmt.Println("Starting to listen to port ", *port)
	e.Logger.Fatal(e.Start(*port))
}
