package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"cleanarch/boiler/internal/user/plugin"
	"cleanarch/boiler/internal/utils/logger"
	// authhttp "cleanarch/boiler/pkg/auth/delivery/http"
	// authmongo "cleanarch/boiler/pkg/auth/repository/mongo"
	// authusecase "cleanarch/boiler/pkg/auth/usecase"
)

type App struct {
	httpServer *http.Server
	userPlugin Plugin
}

type Plugin interface {
	Register()
}

func NewApp() *App {
	err := godotenv.Load("./.env")
	db := initDB()

	r := chi.NewRouter()

	if err != nil {
		log.Fatal("Error loading.env file")
	}
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	r.Use(middleware.Recoverer)

	//compressor := middleware.NewCompressor(5, "text/html", "text/css", "text/plain", "text/javascript", "application/json")
	// compressor.SetEncoder("br", func(w io.Writer, level int) io.Writer {
	// 	params := brotli_enc.NewBrotliParams()
	// 	params.SetQuality(level)
	// 	return brotli_enc.NewBrotliWriter(params, w)
	// })
	// r.Use(compressor.Handler)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	httpServer := &http.Server{
		Addr:           ":" + "3000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger := logger.NewLogger("")
	userPlugin := plugin.NewUserPlugin(r, db, logger)
	userPlugin.Register()

	return &App{
		userPlugin: userPlugin,
		httpServer: httpServer,
	}
}

func (a *App) Run(port string) error {

	// HTTP Server

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatalf("Error occurred while establishing connection to MongoDB: %v", err)
	}
	// It's a good idea to ping the database to verify that the connection is alive
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Ping failed to MongoDB: %v", err)
	}

	return client.Database(os.Getenv("DB_NAME"))
}
