package translateapp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "translateapp/internal/logger"
)

type App struct {
	Service *Service
	Router  *mux.Router
}

type Handler interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
	LanguagePageHandler(writer http.ResponseWriter, request *http.Request)
	TranslatePageHandler(writer http.ResponseWriter, request *http.Request)
	HandleRequests(port string)
	Routes(router *mux.Router)
	GetRouter() *mux.Router
	GetService() *Service
}

func (app *App) GetRouter() *mux.Router {
	return app.Router
}

// Starting Http server on gorilla mux router
func (app *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	app.Router.ServeHTTP(writer, request)
}

// Starting application
func NewApp(service *Service) *App {

	router := mux.NewRouter().StrictSlash(true)

	return &App{
		Service: service,
		Router:  router,
	}
}

// Request to fetch all languages from Libretranslate service.
func (app *App) LanguagePageHandler(writer http.ResponseWriter, request *http.Request) {
	languages, err := app.Service.Languages()
	if err != nil {
		fmt.Fprintf(writer, "Error: %s", err.Error())
	}
	app.Service.Logger.Debug("GET request on localhost:8080/languages")
	app.Service.Logger.Debug("Key: language  Cached value: ", languages)
	if err != nil {
		fmt.Fprintf(writer, "Error: %s", err.Error())
	}
	listOflanguages, _ := json.Marshal(languages)
	fmt.Fprintf(writer, "%s", listOflanguages)
}

// Request to get translation from Libretranslate service.
func (app *App) TranslatePageHandler(writer http.ResponseWriter, request *http.Request) {
	translate, err := app.Service.Translate(request.FormValue("q"), request.FormValue("source"), request.FormValue("target"))
	if err != nil {
		fmt.Fprintf(writer, "Error: %s", err.Error())
	}
	q := request.FormValue("q")

	app.Service.Logger.Debug("Key: "+q, " Cached value: ", translate)
	app.Service.Logger.Debug("POST request on localhost:8080/translate")
	translated, _ := json.Marshal(translate)
	fmt.Fprintf(writer, "%s", translated)
}

// Method to handle all requests
func (app *App) HandleRequests(port string) {
	//create a new router
	router := app.GetRouter()
	app.Routes(router)
	//start and listen to requests
	log.Fatal(http.ListenAndServe(port, router))
}

// Routing,
//if you want to add new route, add it here
func (app *App) Routes(router *mux.Router) {
	router.HandleFunc("/languages", app.LanguagePageHandler).Methods("GET")
	router.HandleFunc("/translate", app.TranslatePageHandler).Methods("POST")
}
