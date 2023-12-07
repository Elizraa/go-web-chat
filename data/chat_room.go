package data

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/elizraa/Globes/config"
	"golang.org/x/crypto/bcrypt"
)

const (
	// PublicRoom is a room open for anyone to join without authentication
	PublicRoom = "public"
	// PrivateRoom is password protected and requires an authentication token in order to process requests
	PrivateRoom = "private"
	// HiddenRoom is a private room that is not listed on public-facing APIs. TODO: Hide this from GET /chats/<id> as well?
	HiddenRoom = "hidden"
)

// Database model
type ChatRoomDB struct {
	Title       string    `json:"title" gorm:"unique;not null"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"visibility"`
	Password    string    `json:"password,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ID          int       `json:"id"`
}

// ChatRoom is a struct representing a chat room
// TODO:  Add Administrator
type ChatRoom struct {
	ChatRoomDB
	Broker  *Broker            `json:"-"`
	Clients map[string]*Client `json:"-"`
}

func GetChatRoomByID(id int) (*ChatRoomDB, error) {
	var chatRoom ChatRoomDB
	if err := config.DB.First(&chatRoom, id).Error; err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func GetAllChatRoom() (*[]ChatRoomDB, error) {
	chatrooms := []ChatRoomDB{}
	err := config.DB.Debug().Model(&ChatRoomDB{}).Limit(100).Find(&chatrooms).Error
	if err != nil {
		return &[]ChatRoomDB{}, err
	}
	return &chatrooms, err
}

func CreateChatRoom(cr *ChatRoomDB) error {
	// validate chat room request
	if apierr, valid := cr.IsValid(); !valid {
		return apierr
	}
	cr.Type = strings.ToLower(cr.Type)
	if cr.Type != PublicRoom {
		pass, err := bcrypt.GenerateFromPassword([]byte(cr.Password), bcrypt.DefaultCost)
		if err != nil {
			return &APIError{
				Code:  104,
				Field: "secret",
			}
		}
		cr.Password = string(pass)
	} else if cr.Type == PublicRoom {
		cr.Password = ""
	}

	cr.CreatedAt = time.Now()
	cr.UpdatedAt = time.Now()
	return config.DB.Create(cr).Error
}

// ToJSON marshals a ChatRoom object in a JSON encoding that can be returned to users
func (cr ChatRoom) ToJSON() (jsonEncoding []byte, err error) {
	// Populate client slice. TODO: Can this be simplified?
	clientsSlice := make([]Client, len(cr.Clients))
	var i int = 0
	for _, v := range cr.Clients {
		//clientsSlice = append(clientsSlice, *v)
		clientsSlice[i] = *v
		i++
	}
	// Create new JSON struct with clients
	jsonEncoding, err = json.Marshal(struct {
		*ChatRoomDB
		Clients []Client `json:"users"`
	}{
		ChatRoomDB: &cr.ChatRoomDB,
		Clients:    clientsSlice,
	})
	return jsonEncoding, err
}

// AddClient will add a user to a ChatRoom
func (cr ChatRoom) AddClient(c *Client) (err error) {
	if cr.clientExists(c.Username) {
		return &APIError{
			Code:  202,
			Field: c.Username,
		}
	}
	cr.Clients[strings.ToLower(c.Username)] = c
	return
}

// RemoveClient will remove a user from a ChatRoom
func (cr ChatRoom) RemoveClient(user string) (err error) {
	if !cr.clientExists(user) {
		return &APIError{
			Code:  201,
			Field: user,
		}
	}
	delete(cr.Clients, strings.ToLower(user))
	return
}

// Authorize authorizes a given ChatEvent for the Room
func (cr ChatRoomDB) Authorize(c *ChatEvent) bool {
	return cr.MatchesPassword(c.Password)
}

// IsValid validates a chat room fields are still valid
func (cr ChatRoomDB) IsValid() (err *APIError, validity bool) {
	// Title should be at least 2 characters
	if len(cr.Title) < 2 || len(cr.Title) > 70 {
		return &APIError{
			Code:  105,
			Field: "title",
		}, false
	}
	// Description shall not be too long
	if len(cr.Description) > 70 {
		return &APIError{
			Code:  105,
			Field: "description",
		}, false
	}
	visibility := strings.ToLower(cr.Type)
	// Visibility must be set
	if visibility != PublicRoom && visibility != PrivateRoom && visibility != HiddenRoom {
		return &APIError{
			Code:  105,
			Field: "visibility",
		}, false
	}
	// Non-public rooms require a valid password
	if (len(cr.Password) < 8) && visibility != PublicRoom {
		return &APIError{
			Code:  105,
			Field: "password",
		}, false
	}
	// A public room should not have a password set (to avoid accidents)
	if len(cr.Password) != 0 && visibility == PublicRoom {
		return &APIError{
			Code:  105,
			Field: "visibility",
		}, false
	}
	return nil, true
}

// MatchesPassword takes in a value and compares it with the room's password
func (cr ChatRoomDB) MatchesPassword(val string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cr.Password), []byte(val))
	return err == nil
}

func (cr ChatRoom) clientExists(name string) bool {
	name = strings.ToLower(name)
	for k := range cr.Clients {
		if k == name {
			return true
		}
	}
	return false
}

// PrettyTime prints the creation date in a pretty format
func (cr ChatRoomDB) PrettyTime() string {
	layout := "Mon Jan _2 15:04"
	return cr.CreatedAt.Format(layout)
}

// Participants prints the # of active clients
func (cr ChatRoom) Participants() int {
	return len(cr.Clients)
}
