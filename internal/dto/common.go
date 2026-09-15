package dto

// ErrorResponse é o formato padrão para respostas de erro da API.
// @Title ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"recurso não encontrado"`
}
