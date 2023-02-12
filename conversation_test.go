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
		conv, _ := NewConversation(user, bot, WithPromptPrefix("I'll ask you some questions."))
		conv.Say("你好")
		conv.Listen("hi")
		conv.Say("你好")
		conv.Listen("hi")
		conv.Say("你好")
		conv.Listen("hi")
		conv.Say("我刚刚问了什么问题")

		record := "test: I am test. You are bot.\n" +
			"I'll ask you some questions.\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 我刚刚问了什么问题\n"
		assert.Equal(t, record, conv.GetRecord())
		prompt := "test: I am test. You are bot.\n" +
			"I'll ask you some questions.\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 你好\n" +
			"bot: hi\n" +
			"test: 我刚刚问了什么问题\n" +
			"bot: "
		assert.Equal(t, prompt, conv.GetPrompt())
	}
	{
		conv, _ := NewConversation(user, bot, WithPromptPrefix("I'll ask you some questions."))
		assert.Equal(t, user, conv.User())
		assert.Equal(t, bot, conv.Bot())
		conv.Say("")
		conv.Listen("")
		assert.Equal(t, 0, conv.history.Length())
	}
}

func TestHistory(t *testing.T) {
	const (
		user = "test"
		bot  = "bot"
	)
	{
		history := &History{}
		history.Add("hello", user)
		history.Add("hi", bot)
		history.Add("how are you?", user)
		history.Add("I'm fine", bot)
		assert.Equal(t, 4, history.Length())
	}
	{
		history := &History{}
		history.Add("", user)
		history.Add("", bot)
		assert.Equal(t, 0, history.Length())
	}
}
