package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	db *gorm.DB
)

type Progress struct {
	gorm.Model
	Task        string
	Description string
}

type Message struct {
	Message string `json:"message"`
}

func main() {
	for {
		_db, err := gorm.Open("mysql", "root:password@tcp(db:3306)/nippo?charset=utf8&parseTime=true&loc=Local")
		if err == nil {
			db = _db
			defer _db.Close()
			break
		}
		fmt.Println(err)
		time.Sleep(3 * time.Second)
	}

	db.AutoMigrate(&Progress{})

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/progresses", postProgress)
	e.GET("/progresses", getProgresses)

	e.Start(":8080")
}

func postProgress(c echo.Context) error {
	req := &struct {
		Task        string `json:"task"`
		Description string `json:"description"`
	}{}

	err := c.Bind(req)
	if err != nil {
		return err
	}

	p := &Progress{
		Task:        req.Task,
		Description: req.Description,
	}
	db.Create(p)
	return c.NoContent(http.StatusCreated)
}

func getProgresses(c echo.Context) error {
	dayStr := c.QueryParam("day")
	if dayStr != "" {
		day, err := time.Parse("2006-1-2", dayStr)
		if err == nil {
			return c.JSON(http.StatusOK, getDayProgress(day))
		}
		return c.JSON(http.StatusBadRequest, Message{"invalid date str"})
	}
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")
	if startStr != "" && endStr != "" {
		start, serr := time.Parse("2006-1-2", startStr)
		end, eerr := time.Parse("2006-1-2", endStr)
		if serr == nil && eerr == nil {
			if start.Before(end) {
				start, end = end, start
			}

			res := make(map[string]*[]Progress)
			for end.Before(start) {
				res[start.Format("2006-1-2")] = getDayProgress(start)
				start = start.Add(-24 * time.Hour)
			}
			res[end.Format("2006-1-2")] = getDayProgress(end)

			return c.JSON(http.StatusOK, res)
		}

		return c.JSON(http.StatusBadRequest, Message{"invalid date str"})
	}

	return c.JSON(http.StatusOK, getDayProgress(time.Now()))

}

func getDayProgress(day time.Time) *[]Progress {
	dayStr := day.Format("2006-1-2")
	nextDayStr := day.Add(time.Hour * 24).Format("2006-1-2")
	progresses := &[]Progress{}
	db.Where("created_at BETWEEN ? AND ?", dayStr, nextDayStr).Find(progresses)
	return progresses
}
