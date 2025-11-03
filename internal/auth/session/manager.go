package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"shanraq.com/internal/auth"
)

const defaultCookieName = "shanraq_session"

// Manager stores authenticated identities in memory and exposes cookie helpers.
type Manager struct {
	mu         sync.RWMutex
	store      map[string]sessionEntry
	ttl        time.Duration
	cookieName string
}

type sessionEntry struct {
	identity auth.Identity
	expires  time.Time
}

// NewManager builds a new in-memory session manager.
func NewManager(ttl time.Duration, cookieName string) *Manager {
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}
	if cookieName == "" {
		cookieName = defaultCookieName
	}
	return &Manager{
		store:      make(map[string]sessionEntry),
		ttl:        ttl,
		cookieName: cookieName,
	}
}

// Create issues a session for the supplied identity and writes a cookie.
func (m *Manager) Create(w http.ResponseWriter, identity auth.Identity) (string, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}

	m.mu.Lock()
	m.store[token] = sessionEntry{identity: identity, expires: time.Now().Add(m.ttl)}
	m.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(m.ttl),
	})

	return token, nil
}

// Destroy removes a session and clears the cookie on the response.
func (m *Manager) Destroy(w http.ResponseWriter, r *http.Request) {
	if token := m.readToken(r); token != "" {
		m.mu.Lock()
		delete(m.store, token)
		m.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

// Identity returns the identity linked to the current request, when one exists.
func (m *Manager) Identity(r *http.Request) (auth.Identity, bool) {
	token := m.readToken(r)
	if token == "" {
		return auth.Identity{}, false
	}

	m.mu.RLock()
	entry, ok := m.store[token]
	m.mu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		if ok {
			m.mu.Lock()
			delete(m.store, token)
			m.mu.Unlock()
		}
		return auth.Identity{}, false
	}
	return entry.identity, true
}

func (m *Manager) readToken(r *http.Request) string {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
