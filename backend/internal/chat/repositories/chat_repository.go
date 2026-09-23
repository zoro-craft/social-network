package repositories

// ChatRepository is the only layer that reads and writes the chat tables.
type ChatRepository struct{}

// NewChatRepository builds a chat repository.
func NewChatRepository() ChatRepository {
	return ChatRepository{}
}