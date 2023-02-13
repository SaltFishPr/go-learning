package main

import (
	"fmt"
	"strings"
)

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type RecordI interface {
	Add(message string, role string)
	Count() int
	Messages() []*Message
	MessagesWithDepth(depth int) []*Message
}

type Conversation struct {
	user         string
	bot          string
	promptPrefix string
	record       RecordI
}

type ConversationOption func(*Conversation)

func NewConversation(user string, bot string, opts ...ConversationOption) (*Conversation, error) {
	if user == "" {
		return nil, fmt.Errorf("user name is required")
	}
	if bot == "" {
		return nil, fmt.Errorf("bot name is required")
	}
	conv := &Conversation{
		user:         user,
		bot:          bot,
		promptPrefix: fmt.Sprintf("%s: I am %s. You are %s.", user, user, bot),
		record:       &Record{},
	}
	for _, opt := range opts {
		opt(conv)
	}
	return conv, nil
}

func WithTopic(topics ...string) ConversationOption {
	return func(c *Conversation) {
		var builder strings.Builder
		for _, topic := range topics {
			builder.WriteByte('\n')
			builder.WriteString(topic)
		}
		c.promptPrefix += builder.String()
	}
}

func (c *Conversation) User() string {
	return c.user
}

func (c *Conversation) Bot() string {
	return c.bot
}

func (c *Conversation) Say(message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}
	c.record.Add(message, c.user)
}

func (c *Conversation) Listen(message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}
	c.record.Add(message, c.bot)
}

func (c *Conversation) GetRecord() []*Message {
	return c.record.Messages()
}

func (c *Conversation) GetPrompt(message string) string {
	const maxDepth = 4

	var builder strings.Builder
	builder.WriteString(c.promptPrefix)
	builder.WriteByte('\n')
	messages := c.record.MessagesWithDepth(maxDepth)
	for _, message := range messages {
		builder.WriteString(fmt.Sprintf("%s: %s", message.Sender, message.Content))
		builder.WriteByte('\n')
	}
	builder.WriteString(fmt.Sprintf("%s: %s", c.user, message))
	builder.WriteByte('\n')
	builder.WriteString(c.bot)
	builder.WriteString(": ")
	return builder.String()
}
