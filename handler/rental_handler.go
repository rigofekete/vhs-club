package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/middleware"
	"github.com/rigofekete/vhs-club-mvc/service"
)

type RentalHandler struct {
	rentalService service.RentalService
}

func NewRentalHandler(s service.RentalService) *RentalHandler {
	return &RentalHandler{rentalService: s}
}

func (h *RentalHandler) RegisterRoutes(r *gin.Engine) {
	app := r.Group("/api/rentals")
	app.GET("/", h.GetAllActiveRentals)

	user := r.Group("/api/rentals")
	user.Use(middleware.UserAuth())
	{
		user.POST("/:id", h.CreateRental)
		user.PATCH("/:id", h.ReturnRental)
		user.GET("/", h.GetMyActiveRentals)
	}

	admin := r.Group("/api/rentals")
	admin.Use(middleware.AdminAuth())
	{
		admin.DELETE("/", h.DeleteAllRentals)
	}
}

func (h *RentalHandler) CreateRental(c *gin.Context) {
	tapeID := c.Param("id")
	userPublicID, ok := middleware.GetUserID(c)
	if !ok {
		_ = c.Error(apperror.ErrUserFieldValidation)
		return
	}

	createdRental, err := h.rentalService.RentTape(c.Request.Context(), tapeID, userPublicID.String())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, RentalSingleResponse(createdRental))
}

func (h *RentalHandler) ReturnRental(c *gin.Context) {
	rentalID := c.Param("id")
	publicID, ok := middleware.GetUserID(c)
	if !ok {
		_ = c.Error(apperror.ErrUserFieldValidation)
		return
	}

	if err := h.rentalService.ReturnTape(c.Request.Context(), publicID.String(), rentalID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *RentalHandler) GetAllActiveRentals(c *gin.Context) {
	rentals, err := h.rentalService.GetAllActiveRentals(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, RentalListResponse(rentals))
}

func (h *RentalHandler) DeleteAllRentals(c *gin.Context) {
	if err := h.rentalService.DeleteAllRentals(c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
