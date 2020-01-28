package chat

import (
	"github.com/deissh/osu-api-server/pkg"
	"github.com/deissh/osu-api-server/pkg/entity"
	"github.com/rs/zerolog/log"
	"net/http"
)

func SendMessage(senderId uint, channelId uint, content string, IsAction bool) (*entity.ChatMessage, error) {
	var message entity.ChatMessage

	err := pkg.Db.Get(
		&message,
		`INSERT INTO message (sender_id, channel_id, content, is_action)
				VALUES ($1, $2, $3, $4)
				RETURNING *`,
		senderId,
		channelId,
		content,
		IsAction,
	)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("message not send")
		return nil, pkg.NewHTTPError(http.StatusBadRequest, "channel_message", "Message not send.")
	}

	err = pkg.Db.Get(
		&message,
		`SELECT
			message.id, message.sender_id, message.channel_id,
			message.created_at, message.content, message.is_action,
			json_build_object('id', u.id, 'username', u.username, 'avatar_url', u.avatar_url,
			    'country_code', u.country_code, 'is_active', u.is_active, 'is_bot', u.is_bot,
			    'is_supporter', u.is_supporter, 'is_online', false) as sender
		FROM message
		INNER JOIN users u on message.sender_id = u.id
		WHERE message.id = $1`,
		message.ID,
	)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("message not send")
		return nil, pkg.NewHTTPError(http.StatusBadRequest, "channel_message", "Message not send.")
	}

	return &message, nil
}

func GetMessages(channelId uint) (*[]entity.ChatMessage, error) {
	var messages []entity.ChatMessage

	err := pkg.Db.Select(
		&messages,
		`SELECT
			message.id, message.sender_id, message.channel_id,
			message.created_at, message.content, message.is_action,
			json_build_object('id', u.id, 'username', u.username, 'avatar_url', u.avatar_url,
			    'country_code', u.country_code, 'is_active', u.is_active, 'is_bot', u.is_bot,
			    'is_supporter', u.is_supporter, 'is_online', false) as sender
		FROM message
		INNER JOIN users u on message.sender_id = u.id
		WHERE message.channel_id = $1`,
		channelId,
	)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("message not send")
		return nil, pkg.NewHTTPError(http.StatusBadRequest, "channel_message", "Message not send.")
	}

	return &messages, nil
}

func GetMessagesAll(userId uint, since uint) (*[]entity.ChatMessage, error) {
	messages := make([]entity.ChatMessage, 0)

	err := pkg.Db.Select(
		&messages,
		`SELECT
			message.id, message.sender_id, message.channel_id,
			message.created_at, message.content, message.is_action,
			json_build_object('id', u.id, 'username', u.username, 'avatar_url', u.avatar_url,
			    'country_code', u.country_code, 'is_active', u.is_active, 'is_bot', u.is_bot,
			    'is_supporter', u.is_supporter, 'is_online', false) as sender
		FROM message
		INNER JOIN user_channels uc on message.channel_id = uc.id
		INNER JOIN users u on message.sender_id = u.id
		WHERE message.sender_id = $1 AND message.id >= $2`,
		userId,
		since,
	)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("message not send")
		return nil, pkg.NewHTTPError(http.StatusBadRequest, "channel_message", "Message not send.")
	}

	return &messages, nil
}
