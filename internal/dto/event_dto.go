package dto

import "time"

// CreateEventInput representa os dados para criação de um guia.
// @Title CreateEventInput
type CreateEventInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
	Day         string `json:"day" binding:"required"`

	// Temporário até implementar autenticação/contexto
	UserID uint `json:"user_id" binding:"required,min=1"`
}

// EventIDInput valida o parâmetro :id na URL.
// @Title EventIDInput
type EventIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateEventInput representa os campos opcionais para atualização (PATCH) de guia.
// @Title UpdateEventInput
type UpdateEventInput struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	BannerURL   string    `json:"banner_url" binding:"omitempty,url"`
	Day         time.Time `json:"day" binding:"required"`
}

// EventResponse representa os dados do guia retornados nas respostas HTTP.
// @Title EventResponse
type EventResponse struct {
	EventID     uint      `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BannerURL   string    `json:"banner_url,omitempty"`
	Day         time.Time `json:"day"`

	UserID uint `json:"user_id"`
}
