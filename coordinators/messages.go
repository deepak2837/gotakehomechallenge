package coordinators

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"take-home-challenge/models"
)

type MessagesCoordinator interface {
	Create(message *models.Message) error
	MarkDeleted(deleteWhen string, messageIDs []int) error
    Read(messageID int) (*models.Message, error)
}

type messagesCoordinator struct {
	db *gorm.DB
}

func NewMessagesCoordinator(db *gorm.DB) MessagesCoordinator {
	return &messagesCoordinator{
		db: db,
	}
}

func (m messagesCoordinator) Create(message *models.Message) error {
	if len(message.UserIDs) == 0 {
		return errors.Wrap(&validationError{message: "UserIDs cannot be empty"}, "validation failed")
	}

	if message.DeletedAt != nil && !message.DeletedAt.IsZero() {
		return errors.Wrap(&validationError{message: "DeletedAt cannot be set on message creation"}, "validation failed")
	}

	if err := m.db.Table("client_messages").Create(message).Error; err != nil {
		return errors.Wrap(err, "failed to create message in database")
	}

	return nil
}

func (m messagesCoordinator) MarkDeleted(deleteWhen string, messageIDs []int) error {
	if len(messageIDs) == 0 {
		return errors.Wrap(&validationError{message: "message IDs cannot be empty"}, "validation failed")
	}

	for i, id := range messageIDs {
        if id <= 0 {
            return errors.Wrap(&validationError{message: "invalid message ID"}, "validation failed")
        }
        
        // Use parameter binding to avoid SQL injection and ensure compatibility with the driver
        err := m.db.Exec(`
            UPDATE client_messages
            SET deleted_at = ?,
                deleted_at_order = ?
            WHERE id = ?`, deleteWhen, i, id).Error
		if err != nil {
			return errors.Wrap(err, "failed to mark message as deleted in database")
		}
	}
	return nil
}

func (m messagesCoordinator) Read(messageID int) (*models.Message, error) {
    var message models.Message
    
    err := m.db.Table("client_messages").
        Where("id = ? AND deleted_at IS NULL", messageID).
        First(&message).Error
    
    if err != nil {
        if gorm.IsRecordNotFoundError(err) {
            return nil, errors.Wrap(err, "message not found")
        }
        return nil, errors.Wrap(err, "failed to read message from database")
    }
    
    return &message, nil
}
