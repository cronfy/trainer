package sessionstorage

import (
	"sync"

	tgb "github.com/cronfy/trainer/internal/telegrambot/domain"
)

type SessionStorage struct {
	mu       sync.RWMutex
	sessions map[tgb.ChatID]tgb.Session
}

func New() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[tgb.ChatID]tgb.Session),
	}
}

func (s *SessionStorage) GetOrCreate(chatId tgb.ChatID, newSession tgb.Session) tgb.Session {
	s.mu.RLock()
	sess, ok := s.sessions[chatId]
	s.mu.RUnlock()
	if ok {
		return sess
	}

	s.mu.Lock()
	s.sessions[chatId] = newSession
	s.mu.Unlock()
	return newSession
}

func (s *SessionStorage) Set(chatId tgb.ChatID, sess tgb.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatId] = sess
}
