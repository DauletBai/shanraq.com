package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
)

// Identity represents a verified user profile returned by an auth provider.
type Identity struct {
	Subject     string
	Email       string
	FullName    string
	PictureURL  string
	Provider    string
	AccessToken string
}

// Provider describes the behaviour required to implement login flows.
type Provider interface {
	AuthCodeURL(state string) (string, error)
	Exchange(ctx context.Context, code string) (Identity, error)
}

// ErrNotConfigured indicates that no external provider has been configured.
var ErrNotConfigured = errors.New("auth provider not configured")

// NoopProvider is a placeholder until an actual implementation (Google OAuth, Meta, etc.) is wired in.
type NoopProvider struct{}

// NewNoopProvider returns a provider that always signals lack of configuration.
func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

func (n *NoopProvider) AuthCodeURL(_ string) (string, error) {
	return "", ErrNotConfigured
}

func (n *NoopProvider) Exchange(_ context.Context, _ string) (Identity, error) {
	return Identity{}, ErrNotConfigured
}

var _ Provider = (*NoopProvider)(nil)

// ProviderRegistry maintains the mapping of provider keys to implementations.
type ProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry constructs a registry optionally preloading provider placeholders.
func NewRegistry(names ...string) *ProviderRegistry {
	reg := &ProviderRegistry{
		providers: make(map[string]Provider),
	}
	for _, name := range names {
		reg.Register(name, NewNoopProvider())
	}
	return reg
}

// Register adds or replaces a provider implementation.
func (r *ProviderRegistry) Register(name string, provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name = normalizeName(name)
	r.providers[name] = provider
}

// Get retrieves a provider by name, returning ErrNotConfigured if missing.
func (r *ProviderRegistry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	name = normalizeName(name)
	provider, ok := r.providers[name]
	if !ok {
		return nil, ErrNotConfigured
	}
	return provider, nil
}

// List exposes the registered provider keys in alphabetical order.
func (r *ProviderRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]string, 0, len(r.providers))
	for key := range r.providers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// DemoOAuthProvider simulates an OAuth provider for local development.
type DemoOAuthProvider struct {
	name string
}

// NewDemoOAuthProvider returns a provider that instantly redirects to callback with a synthetic identity.
func NewDemoOAuthProvider(name string) Provider {
	return &DemoOAuthProvider{name: normalizeName(name)}
}

func (d *DemoOAuthProvider) AuthCodeURL(state string) (string, error) {
	values := url.Values{}
	values.Set("code", fmt.Sprintf("demo-%s-user", d.name))
	if state != "" {
		values.Set("state", state)
	}
	return fmt.Sprintf("/auth/%s/callback?%s", d.name, values.Encode()), nil
}

func (d *DemoOAuthProvider) Exchange(_ context.Context, code string) (Identity, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return Identity{}, errors.New("invalid code")
	}
	subject := fmt.Sprintf("demo-%s-%s", d.name, strings.ReplaceAll(code, " ", "-"))
	email := fmt.Sprintf("%s@demo.shanraq.com", strings.ReplaceAll(code, " ", "-"))
	fullName := strings.Title(strings.ReplaceAll(strings.TrimPrefix(code, "demo-"), "-", " "))
	return Identity{
		Subject:     subject,
		Email:       email,
		FullName:    fullName,
		Provider:    d.name,
		AccessToken: fmt.Sprintf("token-%s", code),
	}, nil
}

var _ Provider = (*DemoOAuthProvider)(nil)
