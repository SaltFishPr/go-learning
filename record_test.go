package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	const (
		user = "test"
		bot  = "bot"
	)
	{
		record := NewRecord()
		record.Add("hello", user)
		record.Add("hi", bot)
		record.Add("how are you?", user)
		record.Add("I'm fine", bot)
		assert.Equal(t, 4, record.Count())
		messages := record.Messages()
		assert.Equal(t, 4, len(messages))
		assert.Equal(t, 4, len(record.MessagesWithDepth(6)))
		assert.Equal(t, 2, len(record.MessagesWithDepth(2)))
		assert.Equal(t, []*Message(nil), record.MessagesWithDepth(-1))
		replyMessages := record.MessagesWithParentIDAndDepth(messages[2].ID, 2)
		assert.Equal(t, 2, len(replyMessages))
		assert.Equal(t, "hi", replyMessages[0].Content)
		assert.Equal(t, "how are you?", replyMessages[1].Content)
		assert.Equal(t, []*Message(nil), record.MessagesWithParentIDAndDepth(messages[2].ID, -1))
		assert.Equal(t, []*Message(nil), record.MessagesWithParentIDAndDepth("1111", 2))
		assert.Equal(t, 3, len(record.MessagesWithParentIDAndDepth(messages[2].ID, 10)))
	}
	{
		record := NewRecord()
		record.Add("", user)
		record.Add("", bot)
		assert.Equal(t, 0, record.Count())
	}
}
