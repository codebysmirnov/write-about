package main

import (
	"github.com/codebysmirnov/write-about/app"
	"github.com/codebysmirnov/write-about/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	conf := config.GetConfig()

	a := &app.App{}
	a.Initialize(conf)

	a.Run(":3000")
}
