package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEvent godoc
// @Summary      Criar um novo evento
// @Description  Cria um novo evento com os dados informados no corpo da requisição.
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        event  body      dto.CreateEventInput  true  "Dados do Evento"
// @Success      201        {object}  dto.EventResponse
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	zap.L().Info("[EventHandler] CreateEvent",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CreateEventInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	data, err := time.Parse("02/01/2006", input.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	event := models.Event{
		Title:       input.Title,
		Description: input.Description,
		BannerURL:   input.BannerURL,
		Day:         data,
		UserID:      input.UserID,
	}

	if err := h.eventService.CreateEvent(&event); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.EventResponse{
		EventID:     event.EventID,
		Title:       event.Title,
		Description: event.Description,
		BannerURL:   event.BannerURL,
		Day:         event.Day,
	}

	c.JSON(http.StatusCreated, response)
}

// GetEventById godoc
// @Summary      Buscar evento por ID
// @Tags         Event
// @Produce      json
// @Param        id   path      int  true  "ID do evento" minimum(1)
// @Success      200  {object}  dto.EventResponse
// @Router       /events/{id} [get]
func (h *EventHandler) GetEventById(c *gin.Context) {
	zap.L().Info("[EventHandler] GetEventById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.EventIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	event, err := h.eventService.GetEventById(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.EventResponse{
		EventID:     event.EventID,
		Title:       event.Title,
		Description: event.Description,
		BannerURL:   event.BannerURL,
		Day:         event.Day,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteEventById godoc
// @Summary      Deletar evento por ID
// @Tags         Event
// @Produce      json
// @Param        id   path      int  true  "ID do evento" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /events/{id} [delete]
func (h *EventHandler) DeleteEventById(c *gin.Context) {
	zap.L().Info("[EventHandler] DeleteEventById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.EventIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.eventService.DeleteEventById(input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "evento deletado com sucesso"})
}

// UpdateEvent godoc
// @Summary      Atualizar dados do evento
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "ID do evento" minimum(1)
// @Param        body  body      dto.UpdateEventInput true  "Dados para atualização"
// @Success      200   {object}  dto.EventResponse
// @Router       /events/{id} [patch]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	zap.L().Info("[EventHandler] UpdateEvent",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var uriInput dto.EventIDInput
	var bodyInput dto.UpdateEventInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	eventModel := models.Event{
		Title:       bodyInput.Title,
		Description: bodyInput.Description,
		BannerURL:   bodyInput.BannerURL,
		Day:         bodyInput.Day,
	}

	updatedEvent, err := h.eventService.UpdateEventById(uriInput.ID, &eventModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		if err.Error() == "nenhum dado recebido para atualização" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.EventResponse{
		EventID:     updatedEvent.EventID,
		Title:       updatedEvent.Title,
		Description: updatedEvent.Description,
		BannerURL:   updatedEvent.BannerURL,
		Day:         updatedEvent.Day,
	}

	c.JSON(http.StatusOK, response)
}
