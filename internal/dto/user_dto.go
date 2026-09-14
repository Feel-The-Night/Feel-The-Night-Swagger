package dto

import "time"

// CreateUserInput representa apenas os dados permitidos no POST /users
type CreateUserInput struct {
	Nickname  string `json:"nickname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	DiscordID string `json:"discordID"`
	Password  string `json:"password" binding:"required"`
}

// UserIDInput valida o parâmetro :id genérico na URL para GET, DELETE, etc.
type UserIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateUserInput permite atualizar campos individualmente (todos opcionais)
type UpdateUser struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	DiscordID string `json:"discordID"`
	Password  string `json:"password"`
}

// Extras
type UserErrorResponse struct {
	Error string `json:"error" example:"usuário não encontrado"`
}

// UserResponse representa os dados do usuário retornados nas respostas HTTP (sem dados sensíveis)
type UserResponse struct {
	UserID           uint      `json:"userID"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	DiscordID        string    `json:"discordID"`
	RegistrationDate time.Time `json:"registrationDate"`
	EventCount       int       `json:"eventCount"`
	GuideCount       int       `json:"guideCount"`
	IsModerator      bool      `json:"isModerator"`
	IsVeteran        bool      `json:"isVeteran"`
}
