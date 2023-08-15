package util

func FormatChannelCollectionName(ChannelId string) string {
	return "channel_" + ChannelId
}

func FormatThreadCollectionName(threadId string) string {
	return "thread_" + threadId
}
