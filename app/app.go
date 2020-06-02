package app

import (
	"factory/in_progress/app/handler"
	"factory/in_progress/app/model"
	"factory/in_progress/config"
	"fmt"
	"log"
	"net/http"

	"github.com/codebysmirnov/api/app/auth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

// App - main app struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
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

	// authorization check middleware
	a.Router.Use(auth.LoggingMiddleware)

	// Routing for handling the articles
	a.Post("/user", a.handleRequest(handler.CreateUser))
	a.Get("/users", a.handleRequest(handler.AllUsers))
	a.Get("/user/{id:[0-9]+}", a.handleRequest(handler.GetUserByID))

	a.Get("/api", a.handleRequest(handler.CheckApii))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
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

// RequestHandlerFunction get db, response and request
type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

// add to params db
func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
