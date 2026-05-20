package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/middleware"
	"github.com/rigofekete/vhs-club-mvc/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	user := r.Group("/api/users")
	user.POST("/login", h.UserLogin)
	user.POST("/", h.CreateUser)

	admin := r.Group("/api/users")
	admin.Use(middleware.AdminAuth())
	{
		admin.POST("/batch", h.CreateUserBatch)
		admin.GET("/:id", h.GetUserByID)
		admin.GET("/", h.GetUsers)
		admin.DELETE("/", h.DeleteAllUsers)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser CreateUserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), newUser.ToModel())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, UserSingleResponse(createdUser))
}

func (h *UserHandler) CreateUserBatch(c *gin.Context) {
	var newUsers CreateUserBatchRequest
	if err := c.ShouldBindJSON(&newUsers); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}

	createdUsers, existingCount, err := h.userService.CreateUserBatch(c.Request.Context(), newUsers.ToModels())
	if err != nil {
		_ = c.Error(err)
		return
	}

	batchResponse := UserBatchResponse{
		Users:         UserListResponse(createdUsers),
		AlreadyExists: *existingCount,
	}

	c.JSON(http.StatusCreated, batchResponse)
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var user UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		_ = c.Error(apperror.WrapValidationError(err))
		return
	}

	loggedUser, err := h.userService.UserLogin(c.Request.Context(), user.ToModel())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse(loggedUser))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UserSingleResponse(user))
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UserListResponse(users))
}

func (h *UserHandler) DeleteAllUsers(c *gin.Context) {
	err := h.userService.DeleteAllUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
