package main

import (
	"database/sql"
	"e-shop-modal/internal/config"
	"e-shop-modal/internal/handlers"
	"e-shop-modal/internal/services"
	"e-shop-modal/internal/store"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := config.LoadConfig()
	fmt.Println(config.DBPort)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s", //"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Conectado a MySQL")
	productStore := store.New(db)
	productService := services.New(productStore)
	handlers := handlers.New(productService)

	http.HandleFunc("/products", handlers.HandleProducts)
	http.HandleFunc("/products/", handlers.HandleProductByID)

	http.ListenAndServe(":8080", nil)
}
