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
		record := &Record{}
		record.Add("hello", user)
		record.Add("hi", bot)
		record.Add("how are you?", user)
		record.Add("I'm fine", bot)
		assert.Equal(t, 4, record.Count())
		assert.Equal(t, 4, len(record.Messages()))
		assert.Equal(t, 4, len(record.MessagesWithDepth(6)))
		assert.Equal(t, 2, len(record.MessagesWithDepth(2)))
		assert.Equal(t, 0, len(record.MessagesWithDepth(-1)))
	}
	{
		record := &Record{}
		record.Add("", user)
		record.Add("", bot)
		assert.Equal(t, 0, record.Count())
	}
}
