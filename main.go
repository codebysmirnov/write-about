package main

import (
	"github.com/codebysmirnov/write-about/app"
	"github.com/codebysmirnov/write-about/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic("failed to get application configuration: " + err.Error())
	}

	a := &app.App{}
	a.Initialize(conf)

	a.Run()
}
