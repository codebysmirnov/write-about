package app

import (
	"factory/in_progress/app/auth"
	"factory/in_progress/app/handler"
	"factory/in_progress/app/handler/user"
	"factory/in_progress/app/model"
	"factory/in_progress/config"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

// App - main app struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Routes []handler.Subroute
}

// Initialize the app
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connected")

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()

	a.Handlers()
}

// Handlers sets the all required routers
func (a *App) Handlers() {

	// routing for users
	u := a.Router.PathPrefix("/user").Subrouter()
	// user route can use only authorized users
	u.Use(auth.NewJWT(os.Getenv("SUPER_KEY")).Middleware)
	a.Register(u, user.New(a.DB))

	// routing for guests
	// g := a.Router.PathPrefix("/guest").Subrouter()
}

// Register add subrouter
func (a *App) Register(r *mux.Router, s handler.Subroute) {
	s.Register(r)
	// create links
	a.Routes = append(a.Routes, s)
}

// Run the app on it's router
func (a *App) Run(host string) {
	// Access to headers
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Authorization", "Content-Type"},
		ExposedHeaders: []string{""},
	})

	handler := c.Handler(a.Router)

	log.Fatal(http.ListenAndServe(host, handler))
}
