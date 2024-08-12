package sessionstorage

import (
	"sync"

	"github.com/cronfy/trainer/internal/telegrambot/domain"
)

type SessionStorage struct {
	mu       sync.RWMutex
	sessions map[domain.ChatID]domain.Session
}

func New() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[domain.ChatID]domain.Session),
	}
}

func (s *SessionStorage) GetOrCreate(chatId domain.ChatID) domain.Session {
	s.mu.RLock()
	sess, ok := s.sessions[chatId]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		sess = domain.Session{State: domain.StateHome}
		s.sessions[chatId] = sess
		s.mu.Unlock()
	}
	return sess
}

func (s *SessionStorage) Set(chatId domain.ChatID, sess domain.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatId] = sess
}
