package main

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Record struct {
	head *recordNode
	tail *recordNode
	m    map[string]*recordNode
}

func NewRecord() *Record {
	return &Record{
		m: make(map[string]*recordNode),
	}
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

	r.m[node.id] = node
}

func (r *Record) Count() int {
	return len(r.m)
}

func (r *Record) getNodeByIndex(index int) *recordNode {
	if index < 0 || index >= r.Count() {
		return nil
	}

	var node *recordNode
	if index < r.Count()/2 {
		node = r.head
		for i := 0; i < index; i++ {
			node = node.next
		}
	} else {
		node = r.tail
		for i := r.Count() - 1; i > index; i-- {
			node = node.prev
		}
	}
	return node
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

	idx := r.Count() - depth
	if idx < 0 {
		idx = 0
	}

	node := r.getNodeByIndex(idx)

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

func (r *Record) MessagesWithParentIDAndDepth(parentId string, depth int) []*Message {
	if depth <= 0 {
		return nil
	}

	node, ok := r.m[parentId]
	if !ok {
		return nil
	}
	end := node.next

	for i := 1; i < depth && node.prev != nil; i++ {
		node = node.prev
	}

	var messages []*Message
	for ; node != nil && node != end; node = node.next {
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
