package services

import "errors"

var (
	// ErrNotAllowed rejects a direct message with no accepted follow in either direction.
	ErrNotAllowed = errors.New("not allowed to message this user")
	// ErrNotMember rejects a group message or history read from a non-member.
	ErrNotMember = errors.New("not a member of this group")
	// ErrContentLong rejects content beyond the agreed length limit.
	ErrContentLong = errors.New("content too long")
)

// maxMessageLength is the content limit agreed in the data model.
const maxMessageLength = 5000

// validateContent rejects empty or over-limit content.
func validateContent(content string) error {
	if content == "" || len([]rune(content)) > maxMessageLength {
		return ErrContentLong
	}
	return nil
}

// validateSelfDM rejects a private message sent to the sender's own id.
func validateSelfDM(senderID, refID int) error {
	if senderID == refID {
		return ErrNotAllowed
	}
	return nil
}