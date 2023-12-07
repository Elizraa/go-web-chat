package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Configuration stores config info of server
type Configuration struct {
	Address      string
	RedisURL     string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

// Config captures parsed input from config.json
var Config Configuration

// Mux contains all the HTTP handlers
var (
	Mux *mux.Router
)

// registerHandlers will register all HTTP handlers
func registerHandlers() *mux.Router {
	api := mux.NewRouter()
	// index
	api.HandleFunc("/", logConsole(index))
	// handle static assets by routing requests from /static/ => "public" directory
	staticDir := "/static/"
	api.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(Config.Static))))

	RegisterChatsAPI(api)
	RegisterChatsHTML(api)

	// Chat Sessions (WebSocket)
	// Do not authorize since you can't add headers to WebSockets. We will do authorization when actually receiving chat messages
	api.Handle("/chats/{ID}/ws", errHandler(authorize(webSocketHandler))).Methods(http.MethodGet)

	api.HandleFunc("/favicon.ico", faviconHandler)

	// Error page
	api.HandleFunc("/err", logConsole(handleError)).Methods(http.MethodGet)
	return api
}
func RegisterChatsAPI(api *mux.Router) {
	// REST-API for chat room [JSON]
	api.Handle("/chats", errHandler(handlePost)).Methods(http.MethodPost)
	api.Handle("/chats/events/{ID}", errHandler(handleGetAllChatEvents)).Methods(http.MethodGet)
	api.Handle("/chats/all", errHandler(handleGetAllChatrooms)).Methods(http.MethodGet)
	api.Handle("/chats/{ID}", errHandler(authorize(handleRoom))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	// Check password matches room
	api.Handle("/chats/{ID}/token", errHandler(login)).Methods(http.MethodPost)

	// Check password matches room
	api.Handle("/chats/{ID}/token/renew", errHandler(renewToken)).Methods(http.MethodGet)

}

func RegisterChatsHTML(api *mux.Router) {
	// List all rooms in [HTML]
	api.HandleFunc("/chats", logConsole(listChats)).Methods(http.MethodGet)

	// Load chat box [HTML]
	api.HandleFunc("/chats/{ID}/chatbox", logConsole(chatbox)).Methods(http.MethodGet)

	// Entrance [HTML]
	api.Handle("/chats/{ID}/entrance", errHandler(joinRoom)).Methods(http.MethodGet)
}

func init() {
	loadConfig()
	loadEnvs()
	loadLog()
	// initialize chat server
	// data.CS.Init()
	Mux = registerHandlers()
}

func loadLog() {
	file, err := os.OpenFile("Globes.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file, err = os.OpenFile("../Globes.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", err)
		}
	}

	// Create a MultiWriter that writes to both the file and os.Stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Use the MultiWriter as the output for the logger
	logger = log.New(multiWriter, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		file, err = os.Open("../config.json")
		if err != nil {
			log.Fatalln("Cannot open config file", err)
		}
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

func loadEnvs() {
	if key, ok := os.LookupEnv("SECRET_KEY"); ok {
		secretKey = key
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "image/x-icon")
	//w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "public/img/favicon.ico")
}
