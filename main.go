package main

import (
	"database/sql"
	"e-shop-modal/internal/config"
	"e-shop-modal/internal/handlers"
	"e-shop-modal/internal/middleware"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"e-shop-modal/internal/utils"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	dsn := cfg.DSN

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Conectado a la DB")
	jwtManager := utils.NewJWTManager(string(cfg.JWTSecret))

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, cfg, jwtManager)
	userHandler := handlers.NewUserHandler(userService)

	orderRepository := repositories.NewOrderRepository(db)
	paymentService := services.NewPaymentService(cfg.MPToken, cfg.NotificationURL, productRepository, orderRepository)
	paymentHandler := handlers.NewPaymentHandler(paymentService, cfg.WebhookSecret)

	googleConfig := config.NewGoogleOAuthConfig(cfg.OAuthIDClient, cfg.OAuthSecretClient, cfg.OAuthRedirectURL, []string{"email", "profile"})
	oauthService := services.NewOAuthService(userRepository, googleConfig)
	oauthHandler := handlers.NewOAuthHandler(oauthService, cfg, jwtManager, googleConfig)

	mux := http.NewServeMux()
	// públicas
	server.HandleFunc(mux, "POST /v1/signup", userHandler.HandleSignUp)
	server.HandleFunc(mux, "POST /v1/login", userHandler.HandleLogIn)
	server.HandleFunc(mux, "GET /v1/oauth/", oauthHandler.GoogleLogin)
	server.HandleFunc(mux, "GET /v1/oauth/callback", oauthHandler.GoogleCallback)
	server.HandleFunc(mux, "POST /v1/webhook", paymentHandler.ConfirmWebhook)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	// protegidas
	server.HandleProtected(
		mux, "GET /v1/profile", userHandler.Profile,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "GET /v1/products", productHandler.GetProducts,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "POST /v1/products", productHandler.CreateProduct,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "GET /v1/products/", productHandler.GetProductByID,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "PUT /v1/products/", productHandler.UpdateProduct,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "DELETE /v1/products/", productHandler.DeleteProduct,
		authMiddleware.Authentication,
	)
	server.HandleProtected(
		mux, "POST /v1/payment", paymentHandler.CreateCheckout,
		authMiddleware.Authentication,
	)

	http.ListenAndServe(":8080", mux)
}
