package models

type Thread struct {
	Id            string
	ChannelId     string // the channel this thread belongs to
	RootMessageId string // the message this thread spawned from
}
