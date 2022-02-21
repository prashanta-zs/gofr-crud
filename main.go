package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
	customerHandler "github.com/gofr-crud/handler"
	customerService "github.com/gofr-crud/service"
	"github.com/gofr-crud/store"
)

func main() {
	app := gofr.New()

	app.Server.ValidateHeaders = false

	store := store.New()
	svc := customerService.New(store)
	h := customerHandler.New(svc)

	app.GET("/customers", h.Get)
	app.GET("/customer/{id}", h.GetByID)
	app.POST("/customer", h.Create)
	app.PUT("/customer", h.Update)
	app.DELETE("/customer/{id}", h.Delete)

	app.Server.HTTP.Port = 8080
	fmt.Println("connected...")

	app.Start()
}
