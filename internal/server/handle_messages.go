package server

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"messaging-engine/internal/db/mongo"
	"messaging-engine/internal/models"
	"time"
)

const (
	NewChannelMessage            = "NEW_CHANNEL_MESSAGE"
	NewThreadMessage             = "NEW_THREAD_MESSAGE"
	UpdateChannelMessage         = "UPDATE_CHANNEL_MESSAGE"
	UpdateThreadMessage          = "UPDATE_THREAD_MESSAGE"
	DeleteChannelMessage         = "DELETE_CHANNEL_MESSAGE"
	DeleteThreadMessage          = "DELETE_THREAD_MESSAGE"
	NewChannelMessageReaction    = "NEW_CHANNEL_MESSAGE_REACTION"
	NewThreadMessageReaction     = "NEW_THREAD_MESSAGE_REACTION"
	DeleteChannelMessageReaction = "DELETE_CHANNEL_MESSAGE_REACTION"
	DeleteThreadMessageReaction  = "DELETE_THREAD_MESSAGE_REACTION"
)

func HandleMessage(message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch message.Type {

	case NewChannelMessage:
		var got models.ChannelMessage
		err := mapstructure.Decode(message.Payload["catache_channel_message"], &got)
		if err != nil {
			logrus.Errorf("error when handling NewChannelMessage: mapstructure.Decode: %v", err)
			return err
		}

		got.DateCreated = time.Now().UTC()

		err = mongo.InsertChannelMessage(ctx, got)
		if err != nil {
			logrus.Errorf("error when handling NewMessage: InsertChannelMessage: %v", err)
			return err
		}

	case NewThreadMessage:
		var got models.ThreadMessage
		err := mapstructure.Decode(message.Payload["catache_thread_message"], &got)
		if err != nil {
			logrus.Errorf("error when handling NewThreadMessage: mapstructure.Decode: %v", err)
			return err
		}

		got.DateCreated = time.Now().UTC()

		err = mongo.InsertThreadMessage(ctx, got)
		if err != nil {
			logrus.Errorf("error when handling NewThreadMessage: InsertMessage: %v", err)
			return err
		}

	case UpdateChannelMessage:
		type expected struct {
			NewChannelMessage models.ChannelMessage `mapstructure:"new_catache_channel_message"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf("error when handling UpdateChannelMessage: mapstructure.Decode: %v", err)
			return err
		}

		err = mongo.UpdateChannelMessage(
			ctx,
			got.NewChannelMessage.ChannelId.String(),
			got.NewChannelMessage.MessageId.String(),
			got.NewChannelMessage,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling UpdateMessage: UpdateChannelMessage: %v",
				err,
			)
			return err
		}

	case UpdateThreadMessage:
		type expected struct {
			NewThreadMessage models.ThreadMessage `mapstructure:"new_catache_thread_message"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf("error when handling UpdateThreadMessage: mapstructure.Decode: %v", err)
			return err
		}

		err = mongo.UpdateThreadMessage(
			ctx,
			got.NewThreadMessage.ThreadId.String(),
			got.NewThreadMessage.MessageId.String(),
			got.NewThreadMessage,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling UpdateMessage: UpdateChannelMessage: %v",
				err,
			)
			return err
		}

	case DeleteChannelMessage:
		type expected struct {
			MessageId       string `mapstructure:"message_id"`
			AuthorAccountId string `mapstructure:"author_account_id"`
			ChannelId       string `mapstructure:"channel_id"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf("error when handling DeleteChannelMessage: mapstructure.Decode: %v", err)
			return err
		}

		err = mongo.DeleteChannelMessage(
			ctx,
			got.ChannelId,
			got.MessageId,
			got.AuthorAccountId,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteChannelMessage: DeleteChannelMessage: %v",
				err,
			)
			return err
		}

	case DeleteThreadMessage:
		type expected struct {
			MessageId       string `mapstructure:"message_id"`
			AuthorAccountId string `mapstructure:"author_account_id"`
			ThreadId        string `mapstructure:"thread_id"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf("error when handling DeleteThreadMessage: mapstructure.Decode: %v", err)
			return err
		}

		err = mongo.DeleteThreadMessage(
			ctx,
			got.ThreadId,
			got.MessageId,
			got.AuthorAccountId,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteThreadMessage: DeleteThreadMessage: %v",
				err,
			)
			return err
		}

	case NewChannelMessageReaction:
		type expected struct {
			MessageId string                 `mapstructure:"message_id"`
			ChannelId string                 `mapstructure:"channel_id"`
			Reaction  models.MessageReaction `mapstructure:"reaction"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf(
				"error when handling NewChannelMessageReaction: mapstructure.Decode: %v",
				err,
			)
			return err
		}

		err = mongo.AddReactionToChannelMessage(
			ctx,
			got.ChannelId,
			got.MessageId,
			got.Reaction,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling NewChannelMessageReaction: AddReactionToChannelMessage: %v",
				err,
			)
			return err
		}

	case NewThreadMessageReaction:
		type expected struct {
			MessageId string                 `mapstructure:"message_id"`
			ThreadId  string                 `mapstructure:"thread_id"`
			Reaction  models.MessageReaction `mapstructure:"reaction"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf(
				"error when handling NewThreadMessageReaction: mapstructure.Decode: %v",
				err,
			)
			return err
		}

		err = mongo.AddReactionToThreadMessage(
			ctx,
			got.ThreadId,
			got.MessageId,
			got.Reaction,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling NewThreadMessageReaction: AddReactionToThreadMessage: %v",
				err,
			)
			return err
		}

	case DeleteChannelMessageReaction:
		type expected struct {
			MessageId        string `mapstructure:"message_id"`
			ReactorAccountId string `mapstructure:"reactor_account_id"`
			ChannelId        string `mapstructure:"channel_id"`
			EmojiUnifiedCode string `mapstructure:"emoji_unified_code"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteChannelMessageReaction: mapstructure.Decode: %v",
				err,
			)
			return err
		}

		err = mongo.RemoveReactionFromChannelMessage(
			ctx,
			got.ChannelId,
			got.MessageId,
			got.ReactorAccountId,
			got.EmojiUnifiedCode,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteChannelMessageReaction: RemoveReactionFromChannelMessage: %v",
				err,
			)
			return err
		}

	case DeleteThreadMessageReaction:
		type expected struct {
			MessageId        string `mapstructure:"message_id"`
			ReactorAccountId string `mapstructure:"reactor_account_id"`
			ThreadId         string `mapstructure:"thread_id"`
			EmojiUnifiedCode string `mapstructure:"emoji_unified_code"`
		}
		var got expected
		err := mapstructure.Decode(message.Payload, &got)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteThreadMessageReaction: mapstructure.Decode: %v",
				err,
			)
			return err
		}

		err = mongo.RemoveReactionFromThreadMessage(
			ctx,
			got.ThreadId,
			got.MessageId,
			got.ReactorAccountId,
			got.EmojiUnifiedCode,
		)
		if err != nil {
			logrus.Errorf(
				"error when handling DeleteThreadMessageReaction: RemoveReactionFromThreadMessage: %v",
				err,
			)
			return err
		}

	}

	return nil
}
