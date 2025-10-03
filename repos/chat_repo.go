package repos

import (
	"database/sql"
	"fmt"

	"github.com/XORbit01/jobseeker-backend/models"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Get or create a conversation between two users (order enforced)
func (r *ChatRepository) getOrCreateConversation(userA, userB int) (int, error) {
	if userA > userB {
		userA, userB = userB, userA
	}

	var conversationID int
	err := r.db.QueryRow(`
		SELECT id FROM conversations
		WHERE participant_one_id = $1 AND participant_two_id = $2
	`, userA, userB).Scan(&conversationID)

	if err == sql.ErrNoRows {
		err = r.db.QueryRow(`
			INSERT INTO conversations (participant_one_id, participant_two_id)
			VALUES ($1, $2) RETURNING id
		`, userA, userB).Scan(&conversationID)
	}
	if err != nil {
		return 0, fmt.Errorf("getOrCreateConversation: %w", err)
	}

	return conversationID, nil
}

// SaveMessage inserts a new message and returns its ID
func (r *ChatRepository) SaveMessage(msg models.Message) (int, error) {
	conversationID, err := r.getOrCreateConversation(msg.SenderID, msg.ReceiverID)
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO messages (conversation_id, sender_id, content)
		VALUES ($1, $2, $3) RETURNING id
	`
	var messageID int
	err = r.db.QueryRow(query, conversationID, msg.SenderID, msg.Content).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("could not save message: %w", err)
	}
	return messageID, nil
}

// GetMessageByID retrieves a single message by ID
func (r *ChatRepository) GetMessageByID(id int) (*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, created_at, is_read
		FROM messages WHERE id = $1
	`
	var msg models.Message
	err := r.db.QueryRow(query, id).Scan(
		&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.Timestamp, &msg.Read,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetMessageByID: %w", err)
	}
	return &msg, nil
}

// GetConversationsForUser lists all conversations the user is in
func (r *ChatRepository) GetConversationsForUser(userID int) ([]models.Conversation, error) {
	query := `
		SELECT id, participant_one_id, participant_two_id, created_at
		FROM conversations
		WHERE participant_one_id = $1 OR participant_two_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetConversationsForUser: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(&c.ID, &c.ParticipantOneID, &c.ParticipantTwoID, &c.CreatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, nil
}

// GetMessagesInConversation retrieves all messages from a conversation (if user is a participant)
func (r *ChatRepository) GetMessagesInConversation(conversationID, userID int) ([]models.Message, error) {
	// Verify access
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM conversations
		WHERE id = $1 AND (participant_one_id = $2 OR participant_two_id = $2)
	`, conversationID, userID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("GetMessagesInConversation auth check: %w", err)
	}
	if count == 0 {
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.Query(`
		SELECT id, conversation_id, sender_id, content, created_at, is_read
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("GetMessagesInConversation: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.Timestamp, &msg.Read); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// MarkMessagesAsReadInConversation sets unread messages in a conversation as read
func (r *ChatRepository) MarkMessagesAsReadInConversation(conversationID, userID int) error {
	// Only update messages not sent by current user
	_, err := r.db.Exec(`
		UPDATE messages
		SET is_read = TRUE
		WHERE conversation_id = $1 AND sender_id != $2 AND is_read = FALSE
	`, conversationID, userID)

	if err != nil {
		return fmt.Errorf("MarkMessagesAsReadInConversation: %w", err)
	}
	return nil
}

// GetConversationsWithStatsForUser returns all conversations + message counts for a user
func (r *ChatRepository) GetConversationsWithStatsForUser(userID int) ([]models.ConversationWithStats, int, error) {
	query := `
		SELECT
			c.id,
			c.participant_one_id,
			c.participant_two_id,
			c.created_at,
			COUNT(m.id) AS message_count
		FROM conversations c
		LEFT JOIN messages m ON m.conversation_id = c.id
		WHERE c.participant_one_id = $1 OR c.participant_two_id = $1
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("GetConversationsWithStatsForUser: %w", err)
	}
	defer rows.Close()

	var convs []models.ConversationWithStats
	for rows.Next() {
		var c models.ConversationWithStats
		if err := rows.Scan(&c.ID, &c.ParticipantOneID, &c.ParticipantTwoID, &c.CreatedAt, &c.MessageCount); err != nil {
			return nil, 0, err
		}
		convs = append(convs, c)
	}

	return convs, len(convs), nil
}
