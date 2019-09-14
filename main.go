package main

import (
	"go-game/models"
	_ "go-game/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// os.Setenv("ZONEINFO", "./lib/time/zoneinfo.zip")

	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	beego.Run()
}
