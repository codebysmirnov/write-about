package app

import (
	"encoding/json"
	"fmt"
	"github.com/Delisa-sama/logger"
	"github.com/codebysmirnov/write-about/app/controller"
	"github.com/codebysmirnov/write-about/app/controller/auth"
	"github.com/codebysmirnov/write-about/app/controller/diary"
	"github.com/codebysmirnov/write-about/app/controller/user"
	"github.com/codebysmirnov/write-about/app/middleware"
	jwtauth "github.com/codebysmirnov/write-about/app/middleware/auth/jwt"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"net/http"
	"os"
	"time"
)

// App - main app struct
type App struct {
	addr   string
	Router *mux.Router
	DB     *gorm.DB
	Routes []controller.Controller
}

// Initialize the app
func (a *App) Initialize(config *config.Config) {
	logFile, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open log file to write")
	}
	defer logFile.Close()

	logger.Init(
		logger.Output(logFile),
		logger.Level(logger.INFO), // TODO: get log-level from config
	)

	cfgString, _ := json.Marshal(config)
	logger.Infof("start app with configuration: %s", string(cfgString))

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("db connected")

	a.addr = fmt.Sprintf("%s:%s", config.Host, config.Port)

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()

	jwt := jwtauth.NewJWT(
		jwtauth.SigningKey(config.SigningKey),
		jwtauth.DefaultExpire(time.Minute*15),
	)

	a.Register(
		a.Router.PathPrefix("/auth").Subrouter(),
		auth.New(a.DB, jwt),
	)

	a.Register(
		a.Router.PathPrefix("/user").Subrouter(),
		user.New(a.DB), jwt,
	)

	a.Register(
		a.Router.PathPrefix("/diary").Subrouter(),
		diary.New(a.DB), jwt,
	)
}

// Register a controller with subroute
func (a *App) Register(r *mux.Router, s controller.Controller, m ...middleware.Middleware) {
	for _, mid := range m {
		r.Use(mid.Middleware)
	}
	s.Register(r)
	a.Routes = append(a.Routes, s)
}

// Run the app on it's router
func (a *App) Run() {
	// Access to headers
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Authorization", "Content-Type"},
		ExposedHeaders: []string{""},
	})

	h := c.Handler(a.Router)

	if err := http.ListenAndServe(a.addr, h); err != nil {
		logger.Fatalf("Server crash: %s", err.Error())
	}
}
