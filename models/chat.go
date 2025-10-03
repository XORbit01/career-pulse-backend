package models

import "time"

type Conversation struct {
	ID               int       `json:"id"`
	ParticipantOneID int       `json:"participant_one_id"`
	ParticipantTwoID int       `json:"participant_two_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id,omitempty"` // populated after save
	SenderID       int       `json:"sender_id"`
	ReceiverID     int       `json:"receiver_id,omitempty"` // used only when sending
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp,omitempty"`
	Read           bool      `json:"read,omitempty"`
}

type MessageInput struct {
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

type ConversationWithStats struct {
	Conversation
	MessageCount int `json:"message_count"`
}
