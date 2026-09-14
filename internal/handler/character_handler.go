package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CharacterHandler struct {
	characterService *services.CharacterService
}

func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

// CreateCharacter godoc
// @Summary Criar um novo personagem
// @Description Cria um novo personagem com os dados informados do body
// @Tags Character
// @Accept json
// @Produce json
// @Param character body dto.CreateCharacterInput true "Dados do Personagem"
// @Success 201 {object} models.Character
// @Router /characters [post]
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	zap.L().Info("[CharacterHandler] CreateCharacter",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CreateCharacterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	character := models.Character{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Type:        input.Type,
	}

	if err := h.characterService.CreateCharacter(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, character)
}

// GetCharacterById godoc
// @Summary      Buscar personagem por ID
// @Tags         Character
// @Produce      json
// @Param        id   path      int  true  "ID do personagem" minimum(1)
// @Success      200  {object}  dto.CharacterResponse
// @Router       /characters/{id} [get]
func (h *CharacterHandler) GetCharacterById(c *gin.Context) {
	zap.L().Info("[CharacterHandler] GetCharacterById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CharacterIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.CharacterErrorResponse{Error: "id inválido"})
		return
	}

	character, err := h.characterService.GetCharacterById(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CharacterErrorResponse{Error: "personagem não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CharacterErrorResponse{Error: err.Error()})
		return
	}

	response := dto.CharacterResponse{
		CharacterID: character.CharacterID,
		Name:        character.Name,
		Description: character.Description,
		Type:        character.Type,
		ImageURL:    character.ImageURL,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteCharacterById godoc
// @Summary      Deletar personagem por ID
// @Tags         Character
// @Produce      json
// @Param        id   path      int  true  "ID do personagem" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /characters/{id} [delete]
func (h *CharacterHandler) DeleteCharacterById(c *gin.Context) {
	zap.L().Info("[CharacterHandler] DeleteCharacterById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CharacterIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.CharacterErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.characterService.DeleteCharacterById(input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CharacterErrorResponse{Error: "personagem não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CharacterErrorResponse{Error: err.Error()})
		return
	}

	// Retorna 200 OK com uma mensagem em formato JSON
	c.JSON(http.StatusOK, gin.H{"message": "personagem deletado com sucesso"})
}

// UpdateCharacter godoc
// @Summary      Atualizar dados do personagem
// @Tags         Character
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "ID do personagem" minimum(1)
// @Param        body  body      dto.UpdateCharacter  true  "Dados para atualização"
// @Success      200   {object}  dto.CharacterResponse
// @Router       /characters/{id} [patch]
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	zap.L().Info("[CharacterHandler] UpdateCharacter",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var uriInput dto.CharacterIDInput
	var bodyInput dto.UpdateCharacter

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.CharacterErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.CharacterErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	// Mapeia o DTO para a struct do Model
	characterModel := models.Character{
		Name:        bodyInput.Name,
		Description: bodyInput.Description,
		Type:        bodyInput.Type,
		ImageURL:    bodyInput.ImageURL,
	}

	// Passa o ID e a referência do Model para a Service
	updatedCharacter, err := h.characterService.UpdateCharacterById(uriInput.ID, &characterModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CharacterErrorResponse{Error: "personagem não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CharacterErrorResponse{Error: err.Error()})
		return
	}

	response := dto.CharacterResponse{
		CharacterID: updatedCharacter.CharacterID,
		Name:        updatedCharacter.Name,
		Description: updatedCharacter.Description,
		Type:        updatedCharacter.Type,
		ImageURL:    updatedCharacter.ImageURL,
	}

	c.JSON(http.StatusOK, response)
}
