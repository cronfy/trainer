package telegrambot

import (
	"sync"
)

type SessionStorage struct {
	mu       sync.RWMutex
	sessions map[ChatID]Session
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[ChatID]Session),
	}
}

func (s *SessionStorage) GetOrCreate(chatId ChatID) Session {
	s.mu.RLock()
	sess, ok := s.sessions[chatId]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		sess = Session{State: StateHome}
		s.sessions[chatId] = sess
		s.mu.Unlock()
	}
	return sess
}

func (s *SessionStorage) Set(chatId ChatID, sess Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatId] = sess
}
