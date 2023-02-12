package main

import (
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

type Conversation struct {
	user         string
	bot          string
	promptPrefix string
	history      *History
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
		history:      &History{},
	}
	for _, opt := range opts {
		opt(conv)
	}
	return conv, nil
}

func WithPromptPrefix(promptPrefix ...string) ConversationOption {
	return func(c *Conversation) {
		var builder strings.Builder
		for _, prefix := range promptPrefix {
			builder.WriteByte('\n')
			builder.WriteString(prefix)
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
	c.history.Add(message, c.user)
}

func (c *Conversation) Listen(message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}
	c.history.Add(message, c.bot)
}

func (c *Conversation) GetRecord() string {
	var builder strings.Builder
	builder.WriteString(c.promptPrefix)
	builder.WriteByte('\n')
	for _, message := range c.history.Messages() {
		builder.WriteString(message)
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (c *Conversation) GetPrompt() string {
	const maxDepth = 5

	var builder strings.Builder
	builder.WriteString(c.promptPrefix)
	builder.WriteByte('\n')
	messages := c.history.Messages()
	if len(messages) > maxDepth {
		messages = messages[len(messages)-maxDepth:]
	}
	for _, message := range messages {
		builder.WriteString(message)
		builder.WriteByte('\n')
	}

	builder.WriteString(c.bot)
	builder.WriteString(": ")
	return builder.String()
}

type History struct {
	head   *HistoryNode
	tail   *HistoryNode
	length int
}

func (h *History) Add(message string, role string) {
	if message == "" {
		return
	}

	node := NewHistoryNode(message, role)
	if h.head == nil {
		h.head = node
		h.tail = node
	} else {
		h.tail.next = node
		node.prev = h.tail
		h.tail = node
	}
	h.length++
}

func (h *History) Length() int {
	return h.length
}

func (h *History) Messages() []string {
	var messages []string
	for node := h.head; node != nil; node = node.next {
		messages = append(messages, fmt.Sprintf("%s: %s", node.role, node.message))
	}
	return messages
}

type HistoryNode struct {
	id      string
	role    string
	message string
	prev    *HistoryNode
	next    *HistoryNode
}

func NewHistoryNode(message string, role string) *HistoryNode {
	return &HistoryNode{
		id:      ulid.Make().String(),
		role:    role,
		message: message,
	}
}
