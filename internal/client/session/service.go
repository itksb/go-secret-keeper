package session

import (
	"errors"
	"fmt"
	"sync"
)

var _ ISession = &AppSession{}

// ISession - session interface
type AppSession struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewAppSession - create new session
func NewAppSession() *AppSession {
	return &AppSession{
		data: make(map[string]interface{}),
		mu:   sync.RWMutex{},
	}
}

// GetValue - get value from session
func (a *AppSession) GetValue(key string) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.data[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("key not found: %s", key))
	}
	return v, nil
}

// SetValue - set value to session
func (a *AppSession) SetValue(key string, value interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data[key] = value
	return nil
}
