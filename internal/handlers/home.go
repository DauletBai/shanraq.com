package handlers

import (
	"net/http"

	"shanraq.com/internal/models"
	"shanraq.com/internal/templates"
)

type HomeHandler struct {
	Tmpl templates.Template
}

func (h *HomeHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	data := struct {
		User *models.User
	}{
		User: user,
	}
	h.Tmpl.ExecuteTemplate(w, "base.html", data)
}