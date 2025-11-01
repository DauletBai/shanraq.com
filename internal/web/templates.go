package web

import (
	"html/template"
	"io"
	"os"
	"sync"
	"time"
)

// Renderer wraps parsed HTML templates and exposes render helpers.
type Renderer struct {
	mu        sync.RWMutex
	templates *template.Template
}

// HomePageData captures the dynamic properties injected into the landing page.
type HomePageData struct {
	PageTitle   string
	BrandName   string
	Description string
	PageID      string
	CurrentYear int
}

// NewRenderer parses templates from the web directory.
func NewRenderer() (*Renderer, error) {
	base := template.New("web").Funcs(template.FuncMap{
		"statusColor": statusColor,
	})

	fsys := os.DirFS("web")
	templates, err := base.ParseFS(fsys, "partials/*.html", "home.html")
	if err != nil {
		return nil, err
	}

	return &Renderer{templates: templates}, nil
}

// RenderHome renders the primary landing page.
func (r *Renderer) RenderHome(w io.Writer, data HomePageData) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if data.PageTitle == "" {
		data.PageTitle = "Discover Global Properties | "
	}
	if data.BrandName == "" {
		data.BrandName = "Shanraq"
	}
	if data.Description == "" {
		data.Description = "Search, compare, and manage international real estate listings from a single platform."
	}
	if data.PageID == "" {
		data.PageID = "blog"
	}
	if data.CurrentYear == 0 {
		data.CurrentYear = time.Now().Year()
	}

	return r.templates.ExecuteTemplate(w, "home.html", data)
}

// Unwrap exposes the underlying template for advanced use cases.
func (r *Renderer) Unwrap() *template.Template {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.templates
}

func statusColor(status string) string {
	switch status {
	case "pending", "queued":
		return "warning"
	case "running", "in_progress":
		return "info"
	case "retry":
		return "secondary"
	case "failed", "error":
		return "danger"
	case "done", "completed", "success":
		return "success"
	default:
		return "secondary"
	}
}
