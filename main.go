/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
