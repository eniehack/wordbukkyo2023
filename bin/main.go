package main

import (
	"log"
	"net/http"
	"os"

	"github.com/WORD-COINS/wordbukkyo2/internal/config"
	"github.com/WORD-COINS/wordbukkyo2/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := config.LoadConfigFile(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	h := new(handler.Handler)
	h.ItemDB = config.DataBase.ItemDB
	h.UserDB = config.DataBase.UserDB

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v0/item", h.ListItem)
	r.Post("/api/v0/item", h.AddItem)
	r.Post("/api/v0/order/new", h.NewOrder)

	log.Fatalln(
		http.ListenAndServe(":3333", r),
	)
}
