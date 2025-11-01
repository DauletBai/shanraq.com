package httpserver

import (
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	"shanraq.com/internal/web"
)

// Deps aggregates dependencies required to build the HTTP router.
type Deps struct {
	Logger   zerolog.Logger
	Config   config.Config
	Renderer *web.Renderer
}
