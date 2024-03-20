package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "配置文件路径")
	flag.Parse()

	conf := NewConfig(confFile)
	chatHandler := newChatCompletionHandler(&conf)

	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/v1/chat/completions", chatHandler.HandleCompletions)

	err := e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}
