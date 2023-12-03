package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/elizraa/chitchat/data"
	"github.com/gorilla/mux"
)

// GET /chats
func listChats(w http.ResponseWriter, r *http.Request) {
	// rooms, err := data.CS.Chats()
	chatrooms, err := data.GetAllChatRoom()
	if err != nil {
		Warning("error fetching chatrooms:", err.Error())
		return
	}
	if err != nil {
		errorMessage(w, r, "Cannot retrieve chats")
	} else {
		// to return back to refreshing page:
		//generateHTML(w, &rooms, "layout", "sidebar", "public.header", "list")
		generateHTMLContent(w, &chatrooms, "list")
		return
	}
}

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "index")
}

// GET /chat/<id>/chatbox
// Default page
func chatbox(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["ID"]; ok {
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			Info("erroneous chats API request", r, err)
		}
		generateHTMLContent(w, &cr.ChatRoomDB, "chat")
	}
}

// GET /chats/<id>/entrance
// Default page
func joinRoom(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	queries := mux.Vars(r)
	if ID, ok := queries["ID"]; ok {
		intID, err := strconv.Atoi(ID)
		if err != nil {
			Info("ID not integer", r, err)
			return err
		}
		cr, err := data.GetChatRoomByID(intID)
		if err != nil {
			Info("erroneous chats API request", r, err)
			return err
		}
		crServer, err := data.CS.Retrieve(cr.Title)
		if crServer == nil {
			crServer = &data.ChatRoom{
				ChatRoomDB: *cr,
			}
			data.CS.PushCR(crServer)
		}

		generateHTML(w, (strings.ToLower(cr.Type) == data.PrivateRoom || cr.Type == data.HiddenRoom), "layout", "sidebar", "public.header", "entrance")
	}
	return
}
