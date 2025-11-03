package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shanraq.com/internal/auth"
)

func TestManagerCreateAndIdentity(t *testing.T) {
	manager := NewManager(1*time.Hour, "")

	resp := httptest.NewRecorder()
	identity := auth.Identity{
		Subject:  "user-123",
		Email:    "user@example.com",
		FullName: "Test User",
		Provider: "demo",
	}
	token, err := manager.Create(resp, identity)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if token == "" {
		t.Fatal("Create() returned empty token")
	}

	cookie := cookieFromRecorder(t, resp)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	gotIdentity, ok := manager.Identity(req)
	if !ok {
		t.Fatal("Identity() returned ok=false, want true")
	}
	if gotIdentity.Subject != identity.Subject {
		t.Fatalf("Identity() subject = %q, want %q", gotIdentity.Subject, identity.Subject)
	}
}

func TestManagerDestroy(t *testing.T) {
	manager := NewManager(1*time.Hour, "")

	resp := httptest.NewRecorder()
	identity := auth.Identity{Subject: "destroy-me"}
	if _, err := manager.Create(resp, identity); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	cookie := cookieFromRecorder(t, resp)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	manager.Destroy(httptest.NewRecorder(), req)

	if _, ok := manager.Identity(req); ok {
		t.Fatal("Identity() returned ok=true after Destroy()")
	}
}

func TestManagerIdentityExpiry(t *testing.T) {
	manager := NewManager(10*time.Millisecond, "")

	resp := httptest.NewRecorder()
	identity := auth.Identity{Subject: "expiring-user"}
	if _, err := manager.Create(resp, identity); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	cookie := cookieFromRecorder(t, resp)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	time.Sleep(20 * time.Millisecond)

	if _, ok := manager.Identity(req); ok {
		t.Fatal("Identity() returned ok=true after expiry")
	}
}

func TestMiddlewareInjectsIdentity(t *testing.T) {
	manager := NewManager(1*time.Hour, "")

	resp := httptest.NewRecorder()
	identity := auth.Identity{Subject: "middleware-user"}
	if _, err := manager.Create(resp, identity); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	cookie := cookieFromRecorder(t, resp)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	var captured auth.Identity
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		captured, ok = IdentityFromContext(r.Context())
		if !ok {
			t.Fatal("IdentityFromContext() returned ok=false, want true")
		}
	})

	mw := Middleware(manager)
	mw(next).ServeHTTP(httptest.NewRecorder(), req)

	if captured.Subject != identity.Subject {
		t.Fatalf("captured subject = %q, want %q", captured.Subject, identity.Subject)
	}
}

func cookieFromRecorder(t *testing.T, rr *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	res := rr.Result()
	defer res.Body.Close()

	cookies := res.Cookies()
	if len(cookies) == 0 {
		t.Fatal("response recorder has no cookies")
	}
	return cookies[0]
}
