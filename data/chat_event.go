package data

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/elizraa/Globes/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Subscribe is used to broadcast a message indicating user has joined ChatRoom
	Subscribe = "join"
	// Broadcast is used to broadcast messages to all subscribed users
	Broadcast = "send"
	// Unsubscribe is used to broadcast a message indicating user has left ChatRoom
	Unsubscribe = "leave"
)

// ChatEvent represents a message event in an associated ChatRoom
type ChatEvent struct {
	EventType string    `json:"event_type,omitempty"`
	User      string    `json:"name,omitempty"`
	RoomID    int       `json:"room_id,omitempty"`
	Color     string    `json:"color,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	Password  string    `json:"secret,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Timestamp time.Time `json:"time,omitempty"`
}

// ValidateEvent ensures data is a valid JSON representation of Chat Event and can be parsed as such
func ValidateEvent(data []byte) (ChatEvent, error) {
	var evt ChatEvent

	if err := json.Unmarshal(data, &evt); err != nil {
		return evt, &APIError{Code: 303}
	}

	if evt.User == "" {
		return evt, &APIError{Code: 303, Field: "name"}
	} else if evt.Msg == "" && strings.ToLower(evt.EventType) == Broadcast {
		return evt, &APIError{Code: 303, Field: "msg"}
	}
	return evt, nil
}

// InsertChatEvent inserts a ChatEvent into MongoDB
func InsertChatEvent(event ChatEvent) error {
	chatEventCollection := config.MongoDBClient.Database("chat_server").Collection("chat_events")

	// Set timestamp if it is not already set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	_, err := chatEventCollection.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}

	return nil
}

// FetchChatEvents retrieves chat events based on roomID and sorts by timestamp in ascending order
func FetchChatEvents(roomID int) ([]ChatEvent, error) {
	chatEventCollection := config.MongoDBClient.Database("chat_server").Collection("chat_events")
	filter := bson.D{{Key: "roomid", Value: roomID}}

	options := options.Find().SetSort(bson.D{{Key: "time", Value: 1}}).SetLimit(50)
	cursor, err := chatEventCollection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Println("Error fetching chat events:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var chatEvents []ChatEvent
	for cursor.Next(context.TODO()) {
		var event ChatEvent
		err := cursor.Decode(&event)
		if err != nil {
			log.Println("Error decoding chat event:", err)
			return nil, err
		}
		chatEvents = append(chatEvents, event)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return chatEvents, nil
}
