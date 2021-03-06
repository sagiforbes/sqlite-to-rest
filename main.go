package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/sagiforbes/sqlite-to-rest/utils"
)

var dbFile = flag.String("f", "", "sqlite db file")
var port = flag.String("p", ":4080", "port to listen to. default 4080")

func doGetTable(c echo.Context) (*utils.QueryResult, error) {
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

	var res *utils.QueryResult
	res, err = utils.DbQuery(*dbFile, sql)
	return res, nil
}

func getTable(c echo.Context) error {
	data, err := doGetTable(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, data)
}

func jsonGetTable(c echo.Context) error {
	data, err := doGetTable(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	res := struct {
		ColumnTypes []string
		Data        []map[string]interface{}
	}{
		ColumnTypes: data.ColumnTypes,
		Data:        make([]map[string]interface{}, len(data.Data)),
	}

	var rec map[string]interface{}
	for rowIdx, row := range data.Data {
		rec = make(map[string]interface{})
		for colIdx, colName := range data.Columns {
			rec[colName] = row[colIdx]
		}
		res.Data[rowIdx] = rec
	}

	return c.JSON(http.StatusOK, res)
}
func csvGetTable(c echo.Context) error {
	srcData, err := doGetTable(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	fileName := fmt.Sprintf("./data_%d.csv", time.Now().UnixNano())
	file, err := os.Create(fileName)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	defer func() {

		os.Remove(fileName)
	}()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write(srcData.Columns)
	var csvRecord []string
	for _, srcRec := range srcData.Data {
		csvRecord = make([]string, len(srcRec))
		for fldIdx, fld := range srcRec {
			csvRecord[fldIdx] = fmt.Sprint(fld)
		}
		csvWriter.Write(csvRecord)
	}

	csvWriter.Flush()
	file.Close()

	return c.File(fileName)

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

	//------------ regular type responses
	e.GET("/:table", getTable)
	e.GET("/:table/count", countTableRecord)

	//------------ json type respons
	g := e.Group("/json")
	g.GET("/:table", jsonGetTable)
	//------------ csv type respons
	csvG := e.Group("/csv")
	csvG.GET("/:table", csvGetTable)

	fmt.Println("Starting to listen to port ", *port)
	e.Logger.Fatal(e.Start(*port))
}
