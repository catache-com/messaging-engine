package mongo

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"messaging-engine/internal/models"
	"messaging-engine/internal/util"
	"time"
)

var MongodbClient mongo.Client

func NewChannel(ctx context.Context, channel models.Channel) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection("channels")

	_, err := messageCollection.InsertOne(ctx, channel)
	return err
}

func NewThread(ctx context.Context, thread models.Thread) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection("threads")

	_, err := messageCollection.InsertOne(ctx, thread)
	return err
}

func InsertChannelMessage(ctx context.Context, message models.ChannelMessage) error {
	catacheDatabase := MongodbClient.Database("catache")

	// insert into general channel messages collection
	channelMessagesCollection := catacheDatabase.Collection(
		util.FormatChannelCollectionName(message.ChannelId.String()),
	)
	_, err := channelMessagesCollection.InsertOne(ctx, message)

	return err
}

func InsertThreadMessage(ctx context.Context, message models.ThreadMessage) error {
	catacheDatabase := MongodbClient.Database("catache")

	// insert into general channel messages collection
	threadMessagesCollection := catacheDatabase.Collection(
		util.FormatThreadCollectionName(message.ThreadId.String()),
	)
	_, err := threadMessagesCollection.InsertOne(ctx, message)

	return err
}

func FindChannelMessagesByChannelId(
	ctx context.Context,
	ChannelId string,
	datetime time.Time,
	pagination int64,
) ([]models.ChannelMessage, error) {
	catacheDatabase := MongodbClient.Database("catache")
	channelMessagesCollection := catacheDatabase.Collection(
		util.FormatChannelCollectionName(ChannelId),
	)

	filter := bson.M{
		"date_created": bson.M{"$lt": datetime},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date_created", -1}})
	findOptions.SetLimit(pagination)

	cursor, err := channelMessagesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %v", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logrus.Errorf("failed to close cursor: %v", err)
		}
	}()

	var ChannelMessages []models.ChannelMessage
	if err := cursor.All(ctx, &ChannelMessages); err != nil {
		return nil, fmt.Errorf(
			"failed to decode ChannelMessages: %v",
			err,
		)
	}

	return ChannelMessages, nil
}

func FindThreadMessagesByThreadId(
	ctx context.Context,
	threadId string,
	datetime time.Time,
	pagination int64,
) ([]models.ThreadMessage, error) {
	catacheDatabase := MongodbClient.Database("catache")
	threadMessagesCollection := catacheDatabase.Collection(
		util.FormatThreadCollectionName(threadId),
	)

	filter := bson.M{
		"date_created": bson.M{"$lt": datetime},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date_created", -1}})
	findOptions.SetLimit(pagination)

	cursor, err := threadMessagesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %v", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logrus.Errorf("failed to close cursor: %v", err)
		}
	}()

	var ThreadMessages []models.ThreadMessage
	if err := cursor.All(ctx, &ThreadMessages); err != nil {
		return nil, fmt.Errorf(
			"failed to decode ThreadMessages: %v",
			err,
		)
	}

	return ThreadMessages, nil
}

func UpdateChannelMessage(
	ctx context.Context,
	ChannelId, MessageId string,
	message models.ChannelMessage,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatChannelCollectionName(ChannelId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{"$set": message}
	_, err := messageCollection.UpdateOne(ctx, filter, update)
	return err
}

func UpdateThreadMessage(
	ctx context.Context,
	threadId, MessageId string,
	message models.ThreadMessage,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatThreadCollectionName(threadId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{"$set": message}
	_, err := messageCollection.UpdateOne(ctx, filter, update)
	return err
}

func DeleteChannelMessage(ctx context.Context, ChannelId, MessageId, AuthorAccountId string) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatChannelCollectionName(ChannelId))

	filter := bson.M{"message_id": MessageId, "account_id": AuthorAccountId}
	_, err := messageCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}

func DeleteThreadMessage(ctx context.Context, threadId, MessageId, AuthorAccountId string) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatThreadCollectionName(threadId))

	filter := bson.M{"message_id": MessageId, "account_id": AuthorAccountId}
	_, err := messageCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}

func AddReactionToChannelMessage(
	ctx context.Context,
	ChannelId, MessageId string,
	reaction models.MessageReaction,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatChannelCollectionName(ChannelId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{"$push": bson.M{"reactions": reaction}}
	_, err := messageCollection.UpdateOne(ctx, filter, update)
	return err
}

func AddReactionToThreadMessage(
	ctx context.Context,
	threadId, MessageId string,
	reaction models.MessageReaction,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatThreadCollectionName(threadId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{"$push": bson.M{"reactions": reaction}}
	_, err := messageCollection.UpdateOne(ctx, filter, update)
	return err
}

func RemoveReactionFromChannelMessage(
	ctx context.Context,
	ChannelId, MessageId, ReactorAccountId string,
	EmojiUnifiedCode string,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatChannelCollectionName(ChannelId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{
		"$pull": bson.M{
			"reactions": bson.M{
				"account_id":         ReactorAccountId,
				"emoji_unified_code": EmojiUnifiedCode,
			},
		},
	}

	_, err := messageCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(
			"failed to remove reaction from message: %v",
			err,
		)
	}

	return nil
}

func RemoveReactionFromThreadMessage(
	ctx context.Context,
	threadId, MessageId, ReactorAccountId string,
	EmojiUnifiedCode string,
) error {
	catacheDatabase := MongodbClient.Database("catache")
	messageCollection := catacheDatabase.Collection(util.FormatThreadCollectionName(threadId))

	filter := bson.M{"message_id": MessageId}
	update := bson.M{
		"$pull": bson.M{
			"reactions": bson.M{
				"account_id":         ReactorAccountId,
				"emoji_unified_code": EmojiUnifiedCode,
			},
		},
	}

	_, err := messageCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(
			"failed to remove reaction from message: %v",
			err,
		)
	}

	return nil
}
