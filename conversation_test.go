package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversation(t *testing.T) {
	const (
		user = "test"
		bot  = "bot"
	)
	{
		conv, err := NewConversation("", bot)
		assert.Nil(t, conv)
		assert.Error(t, err)
	}
	{
		conv, err := NewConversation(user, "")
		assert.Nil(t, conv)
		assert.Error(t, err)
	}
	{
		conv, _ := NewConversation(user, bot, WithTopic("I'll ask you some questions."))
		conv.Say("你好")
		conv.Listen("hi")
		conv.Say("你好")
		conv.Listen("hi")
		conv.Say("你好")
		conv.Listen("hi")
		prompt := "test: I am test. You are bot.\n" +
			"I'll ask you some questions.\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 我刚刚问了什么问题\n" +
			"bot: "
		assert.Equal(t, prompt, conv.GetPrompt("我刚刚问了什么问题"))

		prompt1 := "test: I am test. You are bot.\n" +
			"I'll ask you some questions.\n" +
			"bot: hi\n" +
			"test: 我刚刚问了什么问题\n" +
			"bot: "
		assert.Equal(t, prompt1, conv.GetPrompt("我刚刚问了什么问题", WithDepth(1)))

		prompt2 := "test: I am test. You are bot.\n" +
			"I'll ask you some questions.\n" +
			"test: 你好\n" +
			"test: 我刚刚问了什么问题\n" +
			"bot: "
		assert.Equal(t, prompt2, conv.GetPrompt("我刚刚问了什么问题", WithParentID(conv.GetRecord()[2].ID), WithDepth(1)))
	}

	{
		conv, _ := NewConversation(user, bot, WithTopic("I'll ask you some questions."))
		assert.Equal(t, user, conv.User())
		assert.Equal(t, bot, conv.Bot())
		conv.Say("")
		conv.Listen("")
		assert.Equal(t, 0, conv.record.Count())
	}
}
