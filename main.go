package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/sagiforbes/sqlite-to-rest/utils"
)

var dbFile = flag.String("f", "", "sqlite db file")
var port = flag.String("p", ":4080", "port to listen to. default 4080")

func getTable(c echo.Context) error {
	startRowStr := c.QueryParam("start")
	lengthStr := c.QueryParam("length")
	orderStr := c.QueryParam("order")

	var offset int64
	var limitRecords int64
	var err error
	offset, err = strconv.ParseInt(startRowStr, 10, 64)
	if err != nil {
		offset = 0
	}

	limitRecords, err = strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		limitRecords = 0
	}

	if limitRecords < 1 {
		limitRecords = -1
	}

	sql := fmt.Sprintf("SELECT * FROM %s", c.Param("table"))

	if orderStr != "" {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, orderStr)
	}

	sql = fmt.Sprintf("%s LIMIT %d , %d", sql, offset, limitRecords)

	res, err := utils.DbQuery(*dbFile, sql)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func countTableRecord(c echo.Context) error {
	res, err := utils.DbCount(*dbFile, c.Param("table"))
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf("{\"count\":%d}", res)))
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

	//------------ json type responses
	e.GET("/:table", getTable)
	e.GET("/:table/count", countTableRecord)

	fmt.Println("Starting to listen to port ", *port)
	e.Logger.Fatal(e.Start(*port))
}
