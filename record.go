package main

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Record struct {
	head   *recordNode
	tail   *recordNode
	length int
}

func (r *Record) Add(message string, role string) {
	if message == "" {
		return
	}

	node := NewRecordNode(message, role)
	if r.head == nil {
		r.head = node
		r.tail = node
	} else {
		r.tail.next = node
		node.prev = r.tail
		r.tail = node
	}
	r.length++
}

func (t *Record) Count() int {
	return t.length
}

func (r *Record) Messages() []*Message {
	var messages []*Message
	for node := r.head; node != nil; node = node.next {
		messages = append(messages, &Message{
			ID:        node.id,
			Sender:    node.role,
			Content:   node.message,
			Timestamp: node.time.UnixMilli(),
		})
	}
	return messages
}

func (r *Record) MessagesWithDepth(depth int) []*Message {
	if depth <= 0 {
		return nil
	}

	var node *recordNode
	if r.length <= depth {
		node = r.head
	} else {
		node = r.tail
		for i := 1; i < depth && node != nil; i++ {
			node = node.prev
		}
	}

	var messages []*Message
	for ; node != nil; node = node.next {
		messages = append(messages, &Message{
			ID:        node.id,
			Sender:    node.role,
			Content:   node.message,
			Timestamp: node.time.UnixMilli(),
		})
	}

	return messages
}

type recordNode struct {
	id      string
	role    string
	message string
	time    time.Time
	prev    *recordNode
	next    *recordNode
}

func NewRecordNode(message string, role string) *recordNode {
	return &recordNode{
		id:      ulid.Make().String(),
		role:    role,
		message: message,
		time:    time.Now(),
	}
}
