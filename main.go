package main

import (
	"context"
	"fmt"
	"inout/infrastructure/persistence"
	httpHandler "inout/interface/http"
	"inout/interface/middleware"
	"inout/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

var (
	httpServer *http.Server
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db, err := persistence.NewRepositories()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	fmt.Println(db.Name())
	fmt.Println("Application start")

	router := gin.New()
	router.Use(gin.Recovery())

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "PUT"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://github.com"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))

	router.Use(cors.Default())

	userRepository := persistence.NewUserRepository(db)
	userActivityRepository := persistence.NewUserActivityRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userActivityRepository)
	userHandler := httpHandler.NewUserHandler(userUsecase)

	router.GET("/", func(ctx *gin.Context) {
		fmt.Println("Inout API")
		ctx.JSON(http.StatusOK, "Inout API")
	})
	router.POST("/login", userHandler.Login)

	api := router.Group("api")
	api.Use(middleware.Auth(userRepository, userActivityRepository))

	api.POST("/logout", userHandler.Logout)
	api.POST("/checkin", userHandler.Checkin)
	api.POST("/checkout", userHandler.Checkout)
	api.GET("/history", userHandler.History)
	api.POST("/create", userHandler.Create)
	api.PUT("/update/:id", userHandler.Update)
	api.DELETE("/delete/:id", userHandler.Delete)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	port := os.Getenv("PORT")
	g.Go(func() error {
		httpServer := &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      router,
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	fmt.Println("Running on port", port)

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}

	err = g.Wait()
	if err != nil {
		log.Printf("server returning an error %v", err)
		os.Exit(2)
	}
}
