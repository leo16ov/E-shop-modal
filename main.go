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
	_ "github.com/lib/pq"
)

func main() {
	config := config.LoadConfig()
	dsn := config.DSN

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Conectado a la DB")
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	orderRepository := repositories.NewOrderRepository(db)
	paymentService := services.NewPaymentService(config.MPToken, productRepository, orderRepository)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()
	// públicas
	server.HandleFunc(mux, "POST /signup", userHandler.HandleSignUp)
	server.HandleFunc(mux, "POST /login", userHandler.HandleLogIn)
	server.HandleFunc(mux, "POST /webhook", paymentHandler.ConfirmWebhook)

	// protegidas
	server.HandleProtected(
		mux, "GET /profile", userHandler.Profile,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "GET /products", productHandler.GetProducts,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "POST /products", productHandler.CreateProduct,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "GET /products/", productHandler.GetProductByID,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "PUT /products/", productHandler.UpdateProduct,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "DELETE /products/", productHandler.DeleteProduct,
		middleware.Authentication,
	)
	server.HandleProtected(
		mux, "POST /payment", paymentHandler.CreateCheckout,
		middleware.Authentication,
	)

	http.ListenAndServe(":8080", mux)
}
