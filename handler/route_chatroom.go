package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/elizraa/chitchat/data"
	"github.com/gorilla/mux"
)

// main handler function
func handleRoom(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	queries := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	if ID, ok := queries["ID"]; ok {
		cr, err := data.CS.Retrieve(ID)
		if err != nil {
			Info("erroneous chats API request", r, err)
			return err
		}
		switch r.Method {
		case "GET":
			err = handleGet(w, r, cr)
			return err
		case "PUT":
			err = (handlePut(w, r, cr, ID))
			return err
		case "DELETE":
			err = handleDelete(w, r, cr, ctxId)
			return err
		}
	} else {
		err = &data.APIError{
			Code: 103,
		}
	}

	return err
}

// Retrieve a chat room
// GET /chat/1
func handleGet(w http.ResponseWriter, r *http.Request, cr *data.ChatRoom) (err error) {
	res, err := cr.ToJSON()
	if err != nil {
		return
	}
	Info("retrieved chat room:", cr.Title)
	if _, err := w.Write(res); err != nil {
		Danger("Error writing", res)
	}
	return
}

// Create a ChatRoom
// POST /chats
func handlePost(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		Danger("Error reading", r, err.Error())
		return
	}

	// create ChatRoom obj
	var cr data.ChatRoomDB
	if err = json.Unmarshal(body, &cr); err != nil {
		Warning("error encountered reading POST:", err.Error())
		return err
	}

	if err = data.CreateChatRoom(&cr); err != nil {
		Warning("error creating chatroom:", err.Error())
		return err
	}
	// if err = data.CS.Add(&cr); err != nil {
	// 	Warning("error encountered adding chat room:", err.Error())
	// 	return err
	// }
	// Retrieve updated object
	// createdChatRoom, err := data.CS.Retrieve(cr.Title)
	if err != nil {
		return err
	}
	// res, _ := createdChatRoom.ToJSON()
	w.WriteHeader(201)
	WriteResponse(w, ctxId, cr)
	// if _, err := w.Write(res); err != nil {
	// 	Danger("Error writing", res)
	// }
	return
}

// Update a room
// PUT /chats/<id>
func handlePut(w http.ResponseWriter, r *http.Request, currentChatRoom *data.ChatRoom, title string) (err error) {
	var cr data.ChatRoom
	len := r.ContentLength
	body := make([]byte, len)
	if _, err := r.Body.Read(body); err != nil {
		Danger("Error reading", r, err.Error())
	}
	if err = json.Unmarshal(body, &cr); err != nil {
		Warning("error encountered updating chat room:", err.Error())
		return
	}
	if err = data.CS.Update(title, &cr); err != nil {
		Warning("error encountered updating chat room:", cr, err.Error())
		return
	}
	// Retrieve updated object
	modifiedChatRoom, err := data.CS.RetrieveID(currentChatRoom.ID)
	if err != nil {
		return err
	}
	Info("updated chat room:", title)
	res, _ := modifiedChatRoom.ToJSON()
	if _, err := w.Write(res); err != nil {
		Danger("Error writing", res)
	}
	return
}

// Delete a room
// DELETE /chat/<id>
func handleDelete(w http.ResponseWriter, r *http.Request, cr *data.ChatRoom, ctxId string) (err error) {
	err = data.CS.Delete(cr)
	if err != nil {
		Warning("error encountered deleting chat room:", err.Error())
		return
	}
	// report on status
	Info("deleted chat room:", cr.Title)
	ReportStatus(w, true, nil, ctxId)
	return
}

func handleGetAllChatrooms(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	w.Header().Set("Content-Type", "application/json")

	chatrooms, err := data.GetAllChatRoom()
	if err != nil {
		Warning("error fetching chatrooms:", err.Error())
		return
	}
	w.WriteHeader(201)
	WriteResponse(w, ctxId, chatrooms)
	return
}
