package conversations

import "sync"

var messageLimit = 100

type Message struct {
	Author string
	Body   string
}

type MessageHistory struct {
	Sentinel int
	Ptr      int
	Limit    int

	full     bool
	messages []*Message
	mu       sync.Mutex
}

func (h *MessageHistory) Init() {
	h.Limit = messageLimit + 1
	h.messages = make([]*Message, h.Limit)
}

func (h *MessageHistory) Incr(n int) int {
	if n < 0 || n >= h.Limit {
		panic("n out of range")
	}
	n += 1
	if n == h.Limit {
		n = 0
	}
	return n
}

func (h *MessageHistory) Empty() bool {
	return h.Ptr == h.Sentinel
}

func (h *MessageHistory) Full() bool {
	return h.full
}

func (h *MessageHistory) Advance() {
	h.Ptr = h.Incr(h.Ptr)
	if h.Ptr == h.Sentinel {
		h.full = true
		h.Sentinel = h.Incr(h.Sentinel)
	}
}

func (h *MessageHistory) Append(m *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Advance()
	h.messages[h.Ptr] = m
}

func (h *MessageHistory) GetLastN(n int) []*Message {
	h.mu.Lock()
	defer h.mu.Unlock()

	var recentN []*Message
	if h.Empty() {
		return recentN
	} else if h.Ptr > h.Sentinel {
		recentN = h.messages[h.Sentinel+1 : h.Ptr+1]
	} else {
		recentN = append(h.messages[h.Sentinel+1:h.Limit], h.messages[:h.Ptr+1]...)
	}

	if len(recentN) > n {
		recentN = recentN[len(recentN)-n:]
	}

	return recentN
}

type SlackHistory map[string]*MessageHistory

func (s *SlackHistory) GetChannel(channel string) *MessageHistory {
	if h, ok := (*s)[channel]; ok {
		return h
	} else {
		newHistory := &MessageHistory{}
		newHistory.Init()
		(*s)[channel] = newHistory
		return newHistory
	}
}
