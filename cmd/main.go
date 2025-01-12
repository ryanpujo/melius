package main

import (
	"github.com/ryanpujo/melius/application"
	"github.com/ryanpujo/melius/database"
	"github.com/ryanpujo/melius/internal/jwttoken"
	"github.com/ryanpujo/melius/internal/route"
	"github.com/ryanpujo/melius/registry"
)

func main() {
	db := database.GetDBConnection()
	defer db.Close()
	registry := registry.NewRegistry(db)
	auth := jwttoken.GetJWTAuth(jwttoken.NewTokenVerif())
	app := application.NewApp(route.SetupRoutes(registry.NewAppControllers(), auth))

	if err := app.Serve(); err != nil {
		panic(err)
	}
}
