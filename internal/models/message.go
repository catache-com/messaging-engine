package models

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	FileName string `json:"file_name" mapstructure:"file_name"`
	FileType string `json:"file_type" mapstructure:"file_type"`
}

// Message this is the messages sending to messaging-engine, not exactly user communicated messages
type Message struct {
	Type    string                 `json:"type"    mapstructure:"type"`
	SendTo  string                 `json:"send_to" mapstructure:"send_to"`
	Payload map[string]interface{} `json:"payload" mapstructure:"payload"`
}

type ChannelMessage struct {
	MessageId        uuid.UUID         `bson:"message_id"         json:"message_id"                   mapstructure:"message_id"`
	AuthorAccountId  uuid.UUID         `bson:"author_account_id"  json:"author_account_id"            mapstructure:"author_account_id"`
	ChannelId        uuid.UUID         `bson:"channel_id"         json:"channel_id"                   mapstructure:"channel_id"`
	DateCreated      time.Time         `bson:"date_created"       json:"date_created"                 mapstructure:"date_created"`
	Content          uuid.UUID         `bson:"content"            json:"content"                      mapstructure:"content"`
	Reactions        []MessageReaction `bson:"reactions"          json:"reactions,omitempty"          mapstructure:"reactions"`
	Files            []File            `bson:"files"              json:"files,omitempty"              mapstructure:"files"`
	AttachedThreadId uuid.UUID         `bson:"attached_thread_id" json:"attached_thread_id,omitempty" mapstructure:"attached_thread_id"`
}

type ThreadMessage struct {
	MessageId       uuid.UUID         `bson:"message_id"        json:"message_id"          mapstructure:"message_id"`
	RootMessageId   uuid.UUID         `bson:"root_message_id"   json:"root_message_id"     mapstructure:"root_message_id"`
	AuthorAccountId uuid.UUID         `bson:"author_account_id" json:"author_account_id"   mapstructure:"author_account_id"`
	ThreadId        uuid.UUID         `bson:"thread_id"         json:"thread_id"           mapstructure:"thread_id"`
	DateCreated     time.Time         `bson:"date_created"      json:"date_created"        mapstructure:"date_created"`
	Content         string            `bson:"content"           json:"content"             mapstructure:"content"`
	Reactions       []MessageReaction `bson:"reactions"         json:"reactions,omitempty" mapstructure:"reactions"`
	Files           []File            `bson:"files"             json:"files,omitempty"     mapstructure:"files"`
}

type MessageReaction struct {
	ReactorAccountId uuid.UUID `bson:"reactor_account_id" json:"reactor_account_id" mapstructure:"reactor_account_id"`
	EmojiUnifiedCode string    `bson:"emoji_unified_code" json:"emoji_unified_code" mapstructure:"emoji_unified_code"`
}
