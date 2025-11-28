package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SyedAsadK/govecDB/internal/api"
	"github.com/SyedAsadK/govecDB/internal/db"
)

const (
	STORENAME = "vector.gob"
	PORT      = ":7777"
)

func main() {
	vs, err := db.Load(STORENAME)
	if err != nil {
		vs = db.NewVectorStore()
		fmt.Println("New Database created!")
	} else {
		fmt.Println("Database Loaded!")
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Save every 5 mins
		for range ticker.C {
			fmt.Println("Snapshotting database...")
			err := vs.Save(STORENAME)
			if err != nil {
				fmt.Printf("Error saving DB: %v\n", err)
			}
		}
	}()
	api := api.Controller{Store: vs}
	http.HandleFunc("/vectors", api.HandleInsert)
	http.HandleFunc("/search", api.HandleSearch)
	fmt.Println("Server started on port ", PORT)
	http.ListenAndServe(PORT, nil)
}
