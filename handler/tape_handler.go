package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/middleware"
	"github.com/rigofekete/vhs-club-mvc/service"
)

type TapeHandler struct {
	tapeService service.TapeService
}

// Dependency Injection: TapeHandler depends on the TapeService abstraction
func NewTapeHandler(s service.TapeService) *TapeHandler {
	return &TapeHandler{tapeService: s}
}

func (h *TapeHandler) RegisterRoutes(r *gin.Engine) {
	user := r.Group("/api/tapes")
	user.Use(middleware.UserAuth())
	{
		user.GET("/", h.GetAllTapes)
		user.GET("/:id", h.GetTapeByID)
	}

	admin := r.Group("/api/tapes")
	admin.Use(middleware.AdminAuth())
	{
		admin.POST("/", h.CreateTape)
		admin.POST("/batch/", h.CreateTapeBatch)
		admin.PATCH("/:id", h.UpdateTape)
		admin.DELETE("/:id", h.DeleteTape)
		admin.DELETE("/", h.DeleteAllTapes)
	}
}

func (h *TapeHandler) CreateTape(c *gin.Context) {
	var newTape CreateTapeRequest
	if err := c.ShouldBindJSON(&newTape); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}
	createdTape, err := h.tapeService.CreateTape(c.Request.Context(), newTape.ToModel())
	if err == nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, TapeSingleResponse(createdTape))
}

func (h *TapeHandler) CreateTapeBatch(c *gin.Context) {
	var newTapes CreateTapeBatchRequest
	if err := c.ShouldBindJSON(&newTapes); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}

	createdTapes, existingCount, err := h.tapeService.CreateTapeBatch(c.Request.Context(), newTapes.ToModels())
	if err != nil {
		_ = c.Error(err)
		return
	}

	batchResponse := TapeBatchResponse{
		Tapes:         TapeListResponse(createdTapes),
		AlreadyExists: *existingCount,
	}

	c.JSON(http.StatusCreated, batchResponse)
}

func (h *TapeHandler) GetAllTapes(c *gin.Context) {
	tapes, err := h.tapeService.GetAllTapes(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, TapeListResponse(tapes))
}

func (h *TapeHandler) GetTapeByID(c *gin.Context) {
	id := c.Param("id")
	tape, err := h.tapeService.GetTapeByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, TapeSingleResponse(tape))
}

func (h *TapeHandler) UpdateTape(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTapeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}

	if !updateValid(&req) {
		_ = c.Error(apperror.ErrTapeUpdateRequest)
		return
	}

	tape, err := h.tapeService.UpdateTape(c.Request.Context(), id, req.ToModel())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusPartialContent, TapeUpdateResponse(tape))
}

func (h *TapeHandler) DeleteTape(c *gin.Context) {
	id := c.Param("id")
	if err := h.tapeService.DeleteTape(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TapeHandler) DeleteAllTapes(c *gin.Context) {
	if err := h.tapeService.DeleteAllTapes(c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Helper for UpdateTape
func updateValid(req *UpdateTapeRequest) bool {
	return (req.Title != nil ||
		req.Director != nil ||
		req.Genre != nil ||
		req.Quantity != nil)
}
