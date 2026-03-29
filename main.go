package main

import (
	"database/sql"
	"e-shop-modal/internal/config"
	"e-shop-modal/internal/handlers"
	"e-shop-modal/internal/middleware"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
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
	/*productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)
	*/
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	// públicas
	server.HandleFunc(mux, "POST /signup", userHandler.HandleSignUp)
	server.HandleFunc(mux, "POST /login", userHandler.HandleLogIn)

	// protegidas

	server.HandleProtected(
		mux,
		"GET /profile",
		userHandler.Profile,
		middleware.JWTMiddleware,
	)
	/*
		http.HandleFunc("/products", productHandler.HandleProducts)
		http.HandleFunc("/products/", productHandler.HandleProductByID)*/

	http.ListenAndServe(":8080", mux)
}
