package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type mLog struct {
	User      string
	Namespace string
	PodName   string
	HostIP    string
	PodIP     string
	StartTime time.Time
	UID       string
}

func writeLog(mLog mLog) {
	curDate := time.Now().UTC().Format("2006/01/02")
	dateDir := filepath.Join("/logs", curDate)

	os.MkdirAll(dateDir, os.ModePerm)

	f, err := os.OpenFile(dateDir+"/pod_creations.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()
	log.SetOutput(f)

	log.Printf("%+v", mLog)
}

// POST /api/log
func postLog(c echo.Context) error {
	var mlog mLog

	if err := c.Bind(&mlog); err != nil {
		fmt.Println("Error binding received data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	go writeLog(mlog)

	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/api/log", postLog)
	e.Logger.Info(e.Start(":8080"))
}
