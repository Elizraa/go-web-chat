package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/elizraa/Globes/data"
	"github.com/gorilla/mux"
)

// Add authorization
// POST /login
func user_login(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	w.Header().Set("Content-Type", "application/json")
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	if _, err := io.ReadFull(r.Body, body); err != nil {
		// Handle the error
		Danger("Error reading request", r)
	}
	var user data.UserDB
	if err := json.Unmarshal(body, &user); err != nil {
		Danger("Error parsing token request", r)
	}

	fmt.Println(user)

	userDB, err := data.GetUserByUsername(strings.ToLower(user.Username))
	if err != nil {
		Info("erroneous chats API request", r, err)
		return err
	}
	if userDB.MatchesPassword(user.Password) {
		// Success! Generate token using secret key concatenated with room's password (length > 32)
		tokenString, err := data.EncodeUserJWT(userDB, secretKey)
		if err != nil {
			return err
		}
		// Success, respond with token in JSON body
		jsonEncoding, _ := json.Marshal(struct {
			Outcome  bool   `json:"status"`
			Username string `json:"name"`
			Color    string `json:"color"`
			Token    string `json:"token"`
		}{
			Outcome:  true,
			Username: userDB.Username,
			Color:    userDB.Color,
			Token:    tokenString,
		})
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(jsonEncoding); err != nil {
			Danger("Error writing", jsonEncoding)
		}

	} else {
		return &data.APIError{
			Code:  304,
			Field: "password",
		}
	}
	return
}

// Refreshes tokens before they expire
func user_renewToken(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	w.Header().Set("Content-Type", "application/json")
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
		if cr.Type == data.PublicRoom {
			// Ignore public room
			ReportStatus(w, true, nil, ctxId)
		} else {
			// Check authorization header
			// Get the JWT string from the header
			tknStr, err := extractJwtToken(r)
			if err != nil {
				return &data.APIError{
					Code:  403,
					Field: "token",
				}
			}
			claim := &data.Claims{}
			if err = data.ParseJWT(tknStr, claim, generateUniqueKey(cr)); err != nil {
				return err
			}
			// Success! Generate token
			tokenStringNew, err := claim.RefreshJWT(generateUniqueKey(cr))
			if err != nil {
				return err
			}
			// Success, respond with token in JSON body
			jsonEncoding, _ := json.Marshal(struct {
				Outcome  bool   `json:"status"`
				Username string `json:"name"`
				RoomID   int    `json:"room_id"`
				Token    string `json:"token"`
			}{
				Outcome:  true,
				Username: claim.Username,
				RoomID:   cr.ID,
				Token:    tokenStringNew,
			})
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write(jsonEncoding); err != nil {
				Danger("Error writing", jsonEncoding)
			}
		}
	}
	return
}
